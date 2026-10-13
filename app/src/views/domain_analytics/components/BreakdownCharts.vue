<script setup lang="ts">
import type { DomainAnalytics } from '@/api/nginx_log'
import { Card, Col, Row } from 'antdv-next'
import { storeToRefs } from 'pinia'
import VueApexchart from 'vue3-apexcharts'
import { useSettingsStore } from '@/pinia'

const props = defineProps<{
  data: DomainAnalytics | null
  loading?: boolean
}>()

const settings = useSettingsStore()
const { theme } = storeToRefs(settings)

function fontColor() {
  return theme.value === 'dark' ? '#a3a3a3' : '#525252'
}

const statusLabels = computed(() => {
  const dist = props.data?.status_code_distribution ?? {}
  return Object.keys(dist).sort((a, b) => Number(a) - Number(b))
})

const statusValues = computed(() => {
  const dist = props.data?.status_code_distribution ?? {}
  return statusLabels.value.map(code => dist[code] ?? 0)
})

const statusColors = computed(() => {
  return statusLabels.value.map(code => {
    if (code.startsWith('2'))
      return '#22c55e'
    if (code.startsWith('3'))
      return '#3b82f6'
    if (code.startsWith('4'))
      return '#f59e0b'
    if (code.startsWith('5'))
      return '#ef4444'
    return '#94a3b8'
  })
})

const statusOptions = computed(() => ({
  chart: {
    type: 'donut',
    background: 'transparent',
  },
  labels: statusLabels.value,
  colors: statusColors.value,
  legend: {
    position: 'bottom',
    labels: { colors: fontColor() },
  },
  dataLabels: { enabled: false },
  title: {
    text: $gettext('Status Codes'),
    style: { fontSize: '14px', color: fontColor(), fontWeight: 600 },
  },
  tooltip: {
    theme: theme.value === 'dark' ? 'dark' : 'light',
  },
  plotOptions: {
    pie: {
      donut: {
        size: '62%',
      },
    },
  },
}))

const extensionLabels = computed(() => (props.data?.extensions ?? []).map(item => item.extension))
const extensionValues = computed(() => (props.data?.extensions ?? []).map(item => item.count))

const extensionColors = [
  '#06b6d4',
  '#8b5cf6',
  '#22c55e',
  '#f59e0b',
  '#ef4444',
  '#3b82f6',
  '#ec4899',
  '#14b8a6',
  '#a855f7',
  '#84cc16',
  '#64748b',
  '#f97316',
]

const extensionOptions = computed(() => ({
  chart: {
    type: 'donut',
    background: 'transparent',
  },
  labels: extensionLabels.value,
  colors: extensionColors,
  legend: {
    position: 'bottom',
    labels: { colors: fontColor() },
  },
  dataLabels: { enabled: false },
  title: {
    text: $gettext('Content-Type Breakdown'),
    style: { fontSize: '14px', color: fontColor(), fontWeight: 600 },
  },
  tooltip: {
    theme: theme.value === 'dark' ? 'dark' : 'light',
  },
  plotOptions: {
    pie: {
      donut: {
        size: '62%',
      },
    },
  },
}))
</script>

<template>
  <Row :gutter="[16, 16]" class="mb-4">
    <Col :xs="24" :md="12">
      <Card size="small" class="domain-panel" :loading="loading">
        <VueApexchart
          v-if="statusValues.length"
          type="donut"
          height="300"
          :options="statusOptions"
          :series="statusValues"
        />
        <div v-else class="h-75 flex items-center justify-center text-gray-400">
          {{ $gettext('No status code data') }}
        </div>
      </Card>
    </Col>
    <Col :xs="24" :md="12">
      <Card size="small" class="domain-panel" :loading="loading">
        <VueApexchart
          v-if="extensionValues.length"
          type="donut"
          height="300"
          :options="extensionOptions"
          :series="extensionValues"
        />
        <div v-else class="h-75 flex items-center justify-center text-gray-400">
          {{ $gettext('No content-type data') }}
        </div>
      </Card>
    </Col>
  </Row>
</template>

<style scoped>
.domain-panel {
  border-radius: 10px;
}
</style>
