<script setup lang="ts">
import type { DomainAnalytics } from '@/api/nginx_log'
import { Card, Col, Row, Statistic } from 'antdv-next'
import { storeToRefs } from 'pinia'
import VueApexchart from 'vue3-apexcharts'
import { bytesToSize } from '@/lib/helper'
import { useSettingsStore } from '@/pinia'

const props = defineProps<{
  data: DomainAnalytics | null
  loading?: boolean
}>()

const settings = useSettingsStore()
const { theme } = storeToRefs(settings)

const sparklinePv = computed(() => {
  const hourly = props.data?.hourly_stats ?? []
  return hourly.slice(-24).map(item => item.pv)
})

const sparkOptions = computed(() => ({
  chart: {
    type: 'area',
    sparkline: { enabled: true },
    animations: { enabled: true, speed: 400 },
  },
  stroke: { curve: 'smooth', width: 2 },
  fill: {
    type: 'gradient',
    gradient: {
      opacityFrom: 0.45,
      opacityTo: 0.05,
    },
  },
  colors: [theme.value === 'dark' ? '#3b82f6' : '#2563eb'],
  tooltip: { enabled: false },
}))

const cards = computed(() => {
  const summary = props.data?.summary
  if (!summary) {
    return [
      { key: 'pv', title: $gettext('Total Requests'), value: '—' },
      { key: 'traffic', title: $gettext('Bandwidth'), value: '—' },
      { key: 'success', title: $gettext('Success Rate'), value: '—' },
      { key: 'uv', title: $gettext('Unique Visitors'), value: '—' },
      { key: 'bots', title: $gettext('Bot Traffic'), value: '—' },
      { key: 'humans', title: $gettext('Human Requests'), value: '—' },
    ]
  }

  return [
    {
      key: 'pv',
      title: $gettext('Total Requests'),
      value: summary.total_pv.toLocaleString(),
      spark: true,
    },
    {
      key: 'traffic',
      title: $gettext('Bandwidth'),
      value: bytesToSize(summary.total_traffic),
    },
    {
      key: 'success',
      title: $gettext('Success Rate'),
      value: `${(props.data?.success_rate ?? 0).toFixed(1)}%`,
    },
    {
      key: 'uv',
      title: $gettext('Unique Visitors'),
      value: summary.total_uv.toLocaleString(),
    },
    {
      key: 'bots',
      title: $gettext('Bot Traffic'),
      value: `${(props.data?.bot_percent ?? 0).toFixed(1)}%`,
    },
    {
      key: 'humans',
      title: $gettext('Human Requests'),
      value: (props.data?.human_requests ?? 0).toLocaleString(),
    },
  ]
})
</script>

<template>
  <Row :gutter="[16, 16]" class="mb-4">
    <Col
      v-for="card in cards"
      :key="card.key"
      :xs="12"
      :sm="12"
      :md="6"
    >
      <Card
        size="small"
        class="domain-metric-card"
        :loading="loading"
      >
        <Statistic
          :title="card.title"
          :value="card.value"
        />
        <div
          v-if="card.spark && sparklinePv.length > 1"
          class="mt-2 h-10"
        >
          <VueApexchart
            type="area"
            height="40"
            :options="sparkOptions"
            :series="[{ data: sparklinePv }]"
          />
        </div>
      </Card>
    </Col>
  </Row>
</template>

<style scoped>
.domain-metric-card {
  border-radius: 10px;
  overflow: hidden;
}
</style>
