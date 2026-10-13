<script setup lang="ts">
import type { HealthItem } from '@/api/shield'
import { ReloadOutlined } from '@antdv-next/icons'
import { App, Button, Card, Table, Tag } from 'antdv-next'
import shieldApi from '@/api/shield'

const { message } = App.useApp()
const loading = ref(false)
const items = ref<HealthItem[]>([])

async function load() {
  loading.value = true
  try {
    const res = await shieldApi.health()
    items.value = Array.isArray(res?.data) ? res.data : []
  }
  catch (error) {
    console.error(error)
    message.error($gettext('Failed to load origin health'))
  }
  finally {
    loading.value = false
  }
}

function statusColor(status: string) {
  if (status === 'online')
    return 'success'
  if (status === 'checking')
    return 'processing'
  if (status === 'offline')
    return 'error'
  return 'warning'
}

const columns = computed(() => [
  { title: $gettext('Site'), dataIndex: 'name', key: 'name' },
  { title: $gettext('URL'), dataIndex: 'url', key: 'url' },
  {
    title: $gettext('Status'),
    key: 'status',
    customRender: ({ record }: { record: HealthItem }) => record.status,
  },
  { title: $gettext('Code'), dataIndex: 'status_code', key: 'status_code' },
  { title: $gettext('Latency (ms)'), dataIndex: 'latency_ms', key: 'latency_ms' },
  { title: $gettext('Error'), dataIndex: 'error', key: 'error' },
])

onMounted(() => {
  void load()
})
</script>

<template>
  <Card :title="$gettext('Origin Health')" variant="borderless">
    <template #extra>
      <Button :loading="loading" @click="load">
        <template #icon>
          <ReloadOutlined />
        </template>
        {{ $gettext('Refresh') }}
      </Button>
    </template>

    <p class="mb-4 opacity-70">
      {{ $gettext('Live origin checks from the Site Navigation health monitor — uptime, latency, and failures for every tracked URL.') }}
    </p>

    <Table
      :loading="loading"
      :data-source="items"
      :columns="columns"
      row-key="url"
      :pagination="{ pageSize: 20 }"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <Tag :color="statusColor(record.status)">
            {{ record.status }}
          </Tag>
        </template>
      </template>
    </Table>
  </Card>
</template>
