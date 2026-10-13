<script setup lang="tsx">
import type { TableColumnType } from 'antdv-next'
import type { DomainAnalytics, DomainBlockedIPStat, DomainLogEntry } from '@/api/nginx_log'
import { Card, Col, Row, Table, Tag } from 'antdv-next'
import dayjs from 'dayjs'

const props = defineProps<{
  data: DomainAnalytics | null
  loading?: boolean
}>()

const securityEvents = computed(() => props.data?.security_events ?? [])
const topBlockedIPs = computed(() => props.data?.top_blocked_ips ?? [])

function statusColor(status: number) {
  if (status === 401)
    return 'warning'
  if (status === 403)
    return 'error'
  return 'default'
}

const eventColumns: TableColumnType<DomainLogEntry>[] = [
  {
    title: $gettext('Timestamp'),
    dataIndex: 'timestamp',
    key: 'timestamp',
    width: 170,
    customRender: ({ text }) => dayjs.unix(Number(text)).format('YYYY-MM-DD HH:mm:ss'),
  },
  {
    title: $gettext('Status'),
    dataIndex: 'status',
    key: 'status',
    width: 90,
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
]

const ipColumns: TableColumnType<DomainBlockedIPStat>[] = [
  {
    title: $gettext('Client IP'),
    dataIndex: 'ip',
    key: 'ip',
  },
  {
    title: $gettext('Hits'),
    dataIndex: 'count',
    key: 'count',
    width: 80,
  },
  {
    title: '401',
    dataIndex: 'status_401',
    key: 'status_401',
    width: 70,
  },
  {
    title: '403',
    dataIndex: 'status_403',
    key: 'status_403',
    width: 70,
  },
]
</script>

<template>
  <Card
    size="small"
    class="domain-panel mb-4"
    :title="$gettext('Security & Bouncer')"
    :loading="loading"
  >
    <p class="mb-4 text-sm text-gray-500 dark:text-gray-400">
      {{ $gettext('401 Unauthorized and 403 Forbidden responses, including blocked keyed requests.') }}
    </p>

    <Row :gutter="[16, 16]">
      <Col :xs="24" :lg="14">
        <div class="mb-2 text-sm font-medium">
          {{ $gettext('Blocked / Unauthorized Requests') }}
        </div>
        <Table
          size="small"
          :row-key="(record, index) => `${record.timestamp}-${record.ip}-${record.path}-${index}`"
          :columns="eventColumns"
          :data-source="securityEvents"
          :pagination="false"
          :scroll="{ x: 720, y: 280 }"
        />
      </Col>
      <Col :xs="24" :lg="10">
        <div class="mb-2 text-sm font-medium">
          {{ $gettext('Top 5 Blocked Client IPs') }}
        </div>
        <Table
          size="small"
          row-key="ip"
          :columns="ipColumns"
          :data-source="topBlockedIPs"
          :pagination="false"
        />
      </Col>
    </Row>
  </Card>
</template>

<style scoped>
.domain-panel {
  border-radius: 10px;
}
</style>
