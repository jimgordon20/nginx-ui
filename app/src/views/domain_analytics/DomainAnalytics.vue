<script setup lang="ts">
import type { DomainAnalytics as DomainAnalyticsData, DomainLogEntry, NginxLogData } from '@/api/nginx_log'
import type { DomainTimePreset } from './timeRange'
import { DeleteOutlined, ReloadOutlined } from '@antdv-next/icons'
import { useIntervalFn } from '@vueuse/core'
import { Alert, App, Button, Card, Empty, Select, Segmented, Space } from 'antdv-next'
import nginxLog from '@/api/nginx_log'
import BreakdownCharts from './components/BreakdownCharts.vue'
import RecentRequestsTable from './components/RecentRequestsTable.vue'
import SecurityBouncerPanel from './components/SecurityBouncerPanel.vue'
import SummaryCards from './components/SummaryCards.vue'
import TrafficCharts from './components/TrafficCharts.vue'
import { resolveDomainTimeRange, timePresetLabel } from './timeRange'

const router = useRouter()
const { message, modal } = App.useApp()

const ALL_DOMAINS = '__all__'
const isIndexingEnabled = ref<boolean | null>(null)
const analyticsLoading = ref(false)
const clearingLogs = ref(false)
const domainData = ref<DomainAnalyticsData | null>(null)
const logOptions = ref<{ label: string, value: string }[]>([])
const selectedLogPath = ref(ALL_DOMAINS)
const timePreset = ref<DomainTimePreset>('24h')
const loadError = ref('')

const activeLogPath = computed(() => {
  return selectedLogPath.value === ALL_DOMAINS ? '' : selectedLogPath.value
})

const canClearLogs = computed(() => !!activeLogPath.value)

const recentEntries = computed<DomainLogEntry[]>(() => domainData.value?.recent_entries ?? [])

const timePresetOptions = computed(() => ([
  { label: timePresetLabel('24h'), value: '24h' },
  { label: timePresetLabel('7d'), value: '7d' },
  { label: timePresetLabel('30d'), value: '30d' },
  { label: timePresetLabel('live'), value: 'live' },
]))

const domainSelectOptions = computed(() => [
  { label: $gettext('All Domains'), value: ALL_DOMAINS },
  ...logOptions.value,
])

function friendlyLogLabel(item: NginxLogData) {
  if (item.name)
    return item.name
  if (!item.path)
    return $gettext('Unknown')
  const normalized = item.path.replace(/\\/g, '/')
  const base = normalized.split('/').pop() || normalized
  return base.replace(/_access\.log$/i, '').replace(/\.log$/i, '') || base
}

async function checkIndexingStatus() {
  try {
    const res = await nginxLog.getAdvancedIndexingStatus()
    isIndexingEnabled.value = res.enabled
  }
  catch (error) {
    console.error('Failed to check indexing status:', error)
    isIndexingEnabled.value = false
  }
}

async function loadLogOptions() {
  try {
    const res = await nginxLog.list({ type: 'access' })
    const logs = Array.isArray(res?.data) ? res.data : []
    logOptions.value = logs
      .filter(item => !!item.path)
      .map(item => ({
        label: friendlyLogLabel(item),
        value: item.path!,
      }))
  }
  catch (error) {
    console.error('Failed to load access logs:', error)
    logOptions.value = []
  }
}

async function loadAnalytics() {
  if (isIndexingEnabled.value !== true)
    return

  analyticsLoading.value = true
  loadError.value = ''
  try {
    const range = resolveDomainTimeRange(timePreset.value)
    domainData.value = await nginxLog.getDomainAnalytics({
      log_path: activeLogPath.value || undefined,
      start_time: range.start,
      end_time: range.end,
    })
  }
  catch (error) {
    console.error('Failed to load domain analytics:', error)
    domainData.value = null
    loadError.value = $gettext('Failed to load domain analytics. Make sure log indexing is enabled and logs are indexed.')
  }
  finally {
    analyticsLoading.value = false
  }
}

async function clearDomainLogs() {
  if (!activeLogPath.value)
    return

  modal.confirm({
    title: $gettext('Clear Domain Logs'),
    content: $gettext('This will truncate the selected domain access log file in place. Nginx can keep writing to the same file. Continue?'),
    okText: $gettext('Clear Logs'),
    okType: 'danger',
    cancelText: $gettext('Cancel'),
    async onOk() {
      clearingLogs.value = true
      try {
        await nginxLog.clearDomainLog(activeLogPath.value)
        message.success($gettext('Domain log cleared'))
        await loadAnalytics()
      }
      catch (error) {
        console.error('Failed to clear domain log:', error)
        message.error($gettext('Failed to clear domain log'))
        throw error
      }
      finally {
        clearingLogs.value = false
      }
    },
  })
}

const { pause: pauseLive, resume: resumeLive } = useIntervalFn(async () => {
  if (timePreset.value !== 'live')
    return
  await loadAnalytics()
}, 4000, { immediate: false })

watch(timePreset, preset => {
  if (preset === 'live')
    resumeLive()
  else
    pauseLive()
})

watch([selectedLogPath, timePreset], () => {
  void loadAnalytics()
})

onMounted(async () => {
  await checkIndexingStatus()
  await loadLogOptions()
  if (isIndexingEnabled.value)
    await loadAnalytics()
  if (timePreset.value === 'live')
    resumeLive()
})

onUnmounted(() => {
  pauseLive()
})
</script>

<template>
  <div class="domain-analytics-page">
    <Card
      :title="$gettext('Domain Analytics')"
      variant="borderless"
      class="mb-4 domain-hero-card"
    >
      <template #extra>
        <Space wrap>
          <Button
            danger
            :disabled="!canClearLogs"
            :loading="clearingLogs"
            :title="canClearLogs
              ? $gettext('Truncate the selected domain log file')
              : $gettext('Select a specific domain log to clear')"
            @click="clearDomainLogs"
          >
            <template #icon>
              <DeleteOutlined />
            </template>
            {{ $gettext('Clear Domain Logs') }}
          </Button>
          <Button
            :loading="analyticsLoading"
            @click="loadAnalytics"
          >
            <template #icon>
              <ReloadOutlined />
            </template>
            {{ $gettext('Refresh') }}
          </Button>
        </Space>
      </template>

      <div class="mb-4 flex flex-wrap items-center gap-3 justify-between">
        <div class="flex flex-wrap items-center gap-3">
          <span class="text-sm text-gray-500 dark:text-gray-400">
            {{ $gettext('Domain') }}
          </span>
          <Select
            v-model:value="selectedLogPath"
            class="min-w-64"
            show-search
            :options="domainSelectOptions"
            :placeholder="$gettext('All Domains')"
            :filter-option="(input, option) =>
              ((option?.label ?? '') as string).toLowerCase().includes(input.toLowerCase())"
          />
        </div>

        <Segmented
          v-model:value="timePreset"
          :options="timePresetOptions"
        />
      </div>

      <Alert
        v-if="isIndexingEnabled === false"
        type="warning"
        show-icon
        class="mb-4"
        :message="$gettext('Advanced log indexing is disabled')"
        :description="$gettext('Domain Analytics needs indexed access logs. Enable advanced indexing from Nginx Log → Log List, then rebuild the index.')"
      >
        <template #action>
          <Button type="primary" @click="router.push('/nginx_log/list')">
            {{ $gettext('Open Log List') }}
          </Button>
        </template>
      </Alert>

      <Alert
        v-else-if="loadError"
        type="error"
        show-icon
        class="mb-4"
        :message="loadError"
      />
    </Card>

    <template v-if="isIndexingEnabled">
      <SummaryCards
        :data="domainData"
        :loading="analyticsLoading"
      />

      <TrafficCharts
        :data="domainData"
        :loading="analyticsLoading"
        :preset="timePreset"
      />

      <BreakdownCharts
        :data="domainData"
        :loading="analyticsLoading"
      />

      <SecurityBouncerPanel
        :data="domainData"
        :loading="analyticsLoading"
      />

      <RecentRequestsTable
        :entries="recentEntries"
        :loading="analyticsLoading"
        :is-live="timePreset === 'live'"
      />
    </template>

    <Card
      v-else-if="isIndexingEnabled === false"
      class="domain-empty-card"
    >
      <Empty
        :description="$gettext('Enable advanced indexing to unlock Cloudflare-style domain traffic insights.')"
      />
    </Card>
  </div>
</template>

<style scoped>
.domain-analytics-page {
  min-height: 100%;
}

.domain-hero-card,
.domain-empty-card {
  border-radius: 12px;
}

.min-w-64 {
  min-width: 16rem;
}
</style>
