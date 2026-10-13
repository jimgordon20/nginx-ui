<script setup lang="ts">
import type { DomainAnalytics } from '@/api/nginx_log'
import type { DomainTimePreset } from '../timeRange'
import { Card, Col, Row } from 'antdv-next'
import { storeToRefs } from 'pinia'
import VueApexchart from 'vue3-apexcharts'
import { useSettingsStore } from '@/pinia'

const props = defineProps<{
  data: DomainAnalytics | null
  loading?: boolean
  preset: DomainTimePreset
}>()

const settings = useSettingsStore()
const { theme } = storeToRefs(settings)

function fontColor() {
  return theme.value === 'dark' ? '#a3a3a3' : '#525252'
}

function gridColor() {
  return theme.value === 'dark' ? 'rgba(255,255,255,0.08)' : 'rgba(0,0,0,0.06)'
}

const useHourly = computed(() => props.preset === '24h' || props.preset === 'live')

const categories = computed(() => {
  if (!props.data)
    return [] as string[]

  if (useHourly.value) {
    return (props.data.hourly_stats ?? [])
      .slice(-24)
      .map(item => `${item.hour}:00`)
  }

  return (props.data.daily_stats ?? []).map(item => {
    const date = new Date(item.date)
    return `${date.getMonth() + 1}/${date.getDate()}`
  })
})

const pvSeries = computed(() => {
  if (!props.data)
    return [] as number[]

  if (useHourly.value)
    return (props.data.hourly_stats ?? []).slice(-24).map(item => item.pv)

  return (props.data.daily_stats ?? []).map(item => item.pv)
})

const uvSeries = computed(() => {
  if (!props.data)
    return [] as number[]

  if (useHourly.value)
    return (props.data.hourly_stats ?? []).slice(-24).map(item => item.uv)

  return (props.data.daily_stats ?? []).map(item => item.uv)
})

const baseOptions = computed(() => ({
  chart: {
    type: 'area',
    height: 280,
    toolbar: { show: false },
    zoom: { enabled: false },
    background: 'transparent',
  },
  dataLabels: { enabled: false },
  stroke: { curve: 'smooth', width: 2 },
  fill: {
    type: 'gradient',
    gradient: {
      shadeIntensity: 0.4,
      opacityFrom: 0.35,
      opacityTo: 0.05,
    },
  },
  grid: {
    borderColor: gridColor(),
    strokeDashArray: 3,
  },
  xaxis: {
    categories: categories.value,
    labels: { style: { colors: fontColor() } },
    axisBorder: { show: false },
    axisTicks: { show: false },
  },
  yaxis: {
    labels: {
      style: { colors: fontColor() },
      formatter: (value: number) => Math.round(value).toLocaleString(),
    },
  },
  legend: {
    labels: { colors: fontColor() },
  },
  tooltip: {
    theme: theme.value === 'dark' ? 'dark' : 'light',
  },
}))

const requestsOptions = computed(() => ({
  ...baseOptions.value,
  colors: ['#3b82f6', '#22c55e'],
  title: {
    text: $gettext('Requests Over Time'),
    style: { fontSize: '14px', color: fontColor(), fontWeight: 600 },
  },
}))

const visitorsOptions = computed(() => ({
  ...baseOptions.value,
  colors: ['#8b5cf6'],
  title: {
    text: $gettext('Unique Visitors Over Time'),
    style: { fontSize: '14px', color: fontColor(), fontWeight: 600 },
  },
}))
</script>

<template>
  <Row :gutter="[16, 16]" class="mb-4">
    <Col :xs="24" :lg="14">
      <Card size="small" class="domain-panel" :loading="loading">
        <VueApexchart
          v-if="pvSeries.length"
          type="area"
          height="280"
          :options="requestsOptions"
          :series="[
            { name: $gettext('Requests'), data: pvSeries },
            { name: $gettext('Unique Visitors'), data: uvSeries },
          ]"
        />
        <div v-else class="h-70 flex items-center justify-center text-gray-400">
          {{ $gettext('No traffic data for this range') }}
        </div>
      </Card>
    </Col>
    <Col :xs="24" :lg="10">
      <Card size="small" class="domain-panel" :loading="loading">
        <VueApexchart
          v-if="uvSeries.length"
          type="area"
          height="280"
          :options="visitorsOptions"
          :series="[{ name: $gettext('Unique Visitors'), data: uvSeries }]"
        />
        <div v-else class="h-70 flex items-center justify-center text-gray-400">
          {{ $gettext('No visitor data for this range') }}
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
