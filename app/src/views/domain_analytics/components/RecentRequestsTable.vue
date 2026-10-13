<script setup lang="tsx">
import type { TableColumnType } from 'antdv-next'
import type { DomainLogEntry } from '@/api/nginx_log'
import { Card, Table, Tag } from 'antdv-next'
import dayjs from 'dayjs'
import { bytesToSize } from '@/lib/helper'

defineProps<{
  entries: DomainLogEntry[]
  loading?: boolean
  isLive?: boolean
}>()

function statusColor(status: number) {
  if (status >= 200 && status < 300)
    return 'success'
  if (status >= 300 && status < 400)
    return 'processing'
  if (status >= 400 && status < 500)
    return 'warning'
  if (status >= 500)
    return 'error'
  return 'default'
}

const columns: TableColumnType<DomainLogEntry>[] = [
  {
    title: $gettext('Timestamp'),
    dataIndex: 'timestamp',
    key: 'timestamp',
    width: 180,
    customRender: ({ text }) => dayjs.unix(Number(text)).format('YYYY-MM-DD HH:mm:ss'),
  },
  {
    title: $gettext('Method'),
    dataIndex: 'method',
    key: 'method',
    width: 90,
  },
  {
    title: $gettext('Path'),
    dataIndex: 'path',
    key: 'path',
    ellipsis: true,
  },
  {
    title: $gettext('Status'),
    dataIndex: 'status',
    key: 'status',
    width: 100,
    customRender: ({ text }) => (
      <Tag color={statusColor(Number(text))}>{text}</Tag>
    ),
  },
  {
    title: $gettext('Client IP'),
    dataIndex: 'ip',
    key: 'ip',
    width: 140,
  },
  {
    title: $gettext('Bytes'),
    dataIndex: 'bytes_sent',
    key: 'bytes_sent',
    width: 110,
    customRender: ({ text }) => bytesToSize(Number(text) || 0),
  },
]
</script>

<template>
  <Card
    size="small"
    class="domain-panel"
    :title="isLive ? $gettext('Live Requests') : $gettext('Recent Requests')"
  >
    <template v-if="isLive" #extra>
      <span class="live-dot" />
      <span class="text-xs text-gray-500 dark:text-gray-400">
        {{ $gettext('Updating every few seconds') }}
      </span>
    </template>

    <Table
      size="small"
      :row-key="(record, index) => `${record.timestamp}-${record.ip}-${record.path}-${index}`"
      :columns="columns"
      :data-source="entries"
      :loading="loading"
      :pagination="false"
      :scroll="{ x: 900 }"
    />
  </Card>
</template>

<style scoped>
.domain-panel {
  border-radius: 10px;
}

.live-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  margin-right: 8px;
  border-radius: 9999px;
  background: #22c55e;
  box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.55);
  animation: live-pulse 1.6s ease-out infinite;
  vertical-align: middle;
}

@keyframes live-pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.55);
  }
  70% {
    box-shadow: 0 0 0 8px rgba(34, 197, 94, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(34, 197, 94, 0);
  }
}
</style>
