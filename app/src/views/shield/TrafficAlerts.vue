<script setup lang="ts">
import type { AlertEvaluation, AlertRule } from '@/api/shield'
import { DeleteOutlined, PlusOutlined, ThunderboltOutlined } from '@antdv-next/icons'
import {
  App,
  Button,
  Card,
  Form,
  Input,
  InputNumber,
  Select,
  Space,
  Switch,
  Table,
  Tag,
} from 'antdv-next'
import shieldApi from '@/api/shield'

const { message } = App.useApp()

const loading = ref(false)
const evaluating = ref(false)
const rules = ref<AlertRule[]>([])
const evaluations = ref<AlertEvaluation[]>([])

const draft = reactive<AlertRule>({
  name: '',
  enabled: true,
  site: '',
  log_path: '',
  metric: 'requests',
  threshold: 10000,
  window_minutes: 60,
  notify_on_trigger: true,
})

const metricOptions = computed(() => [
  { label: $gettext('Requests'), value: 'requests' },
  { label: $gettext('Bandwidth (bytes)'), value: 'bandwidth' },
  { label: $gettext('Blocked (401/403)'), value: 'blocked' },
  { label: $gettext('Server error rate %'), value: 'server_errors' },
  { label: $gettext('Client error rate %'), value: 'client_errors' },
])

async function load() {
  loading.value = true
  try {
    const res = await shieldApi.listAlerts()
    rules.value = Array.isArray(res?.data) ? res.data : []
  }
  catch (error) {
    console.error(error)
    message.error($gettext('Failed to load alerts'))
  }
  finally {
    loading.value = false
  }
}

async function saveRule() {
  if (!draft.name.trim()) {
    message.warning($gettext('Name is required'))
    return
  }
  loading.value = true
  try {
    const res = await shieldApi.saveAlert({ ...draft })
    rules.value = res.data
    draft.name = ''
    message.success($gettext('Alert saved'))
  }
  catch (error) {
    console.error(error)
    message.error($gettext('Failed to save alert'))
  }
  finally {
    loading.value = false
  }
}

async function removeRule(id?: string) {
  if (!id)
    return
  const res = await shieldApi.deleteAlert(id)
  rules.value = res.data
  message.success($gettext('Deleted'))
}

async function runEvaluate() {
  evaluating.value = true
  try {
    const res = await shieldApi.evaluateAlerts()
    evaluations.value = Array.isArray(res?.data) ? res.data : []
    message.success($gettext('Evaluation complete'))
  }
  catch (error) {
    console.error(error)
    message.error($gettext('Evaluation failed — indexing may be disabled'))
  }
  finally {
    evaluating.value = false
  }
}

const columns = computed(() => [
  { title: $gettext('Name'), dataIndex: 'name', key: 'name' },
  { title: $gettext('Metric'), dataIndex: 'metric', key: 'metric' },
  { title: $gettext('Threshold'), dataIndex: 'threshold', key: 'threshold' },
  { title: $gettext('Window (min)'), dataIndex: 'window_minutes', key: 'window_minutes' },
  { title: $gettext('Enabled'), key: 'enabled' },
  { title: $gettext('Actions'), key: 'actions' },
])

onMounted(() => {
  void load()
})
</script>

<template>
  <div>
    <Card
      class="mb-4"
      :title="$gettext('Traffic Alerts')"
      variant="borderless"
    >
      <template #extra>
        <Button type="primary" :loading="evaluating" @click="runEvaluate">
          <template #icon>
            <ThunderboltOutlined />
          </template>
          {{ $gettext('Run Now') }}
        </Button>
      </template>

      <p class="mb-4 opacity-70">
        {{ $gettext('Watch Domain Analytics metrics and notify when request, bandwidth, or block spikes cross your limits.') }}
      </p>

      <Form layout="vertical" class="mb-4">
        <div class="grid gap-3 md:grid-cols-2">
          <Form.Item :label="$gettext('Name')">
            <Input v-model:value="draft.name" />
          </Form.Item>
          <Form.Item :label="$gettext('Metric')">
            <Select v-model:value="draft.metric" :options="metricOptions" />
          </Form.Item>
          <Form.Item :label="$gettext('Threshold')">
            <InputNumber v-model:value="draft.threshold" class="w-full" />
          </Form.Item>
          <Form.Item :label="$gettext('Window (minutes)')">
            <InputNumber v-model:value="draft.window_minutes" class="w-full" />
          </Form.Item>
          <Form.Item :label="$gettext('Log path (optional, empty = all)')">
            <Input v-model:value="draft.log_path" />
          </Form.Item>
          <Form.Item :label="$gettext('Notify')">
            <Switch v-model:checked="draft.notify_on_trigger" />
          </Form.Item>
        </div>
        <Button type="dashed" @click="saveRule">
          <template #icon>
            <PlusOutlined />
          </template>
          {{ $gettext('Add Alert') }}
        </Button>
      </Form>

      <Table
        :loading="loading"
        :data-source="rules"
        :columns="columns"
        row-key="id"
        :pagination="false"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'enabled'">
            <Tag :color="record.enabled ? 'success' : 'default'">
              {{ record.enabled ? $gettext('On') : $gettext('Off') }}
            </Tag>
          </template>
          <template v-else-if="column.key === 'actions'">
            <Button type="text" danger @click="removeRule(record.id)">
              <template #icon>
                <DeleteOutlined />
              </template>
            </Button>
          </template>
        </template>
      </Table>
    </Card>

    <Card
      v-if="evaluations.length"
      :title="$gettext('Last Evaluation')"
      variant="borderless"
    >
      <Space direction="vertical" class="w-full">
        <div
          v-for="item in evaluations"
          :key="item.rule.id || item.rule.name"
          class="flex items-center justify-between gap-3"
        >
          <div>
            <div class="font-medium">
              {{ item.rule.name }}
            </div>
            <div class="text-sm opacity-70">
              {{ item.message }}
            </div>
          </div>
          <Tag :color="item.triggered ? 'error' : 'success'">
            {{ item.triggered ? $gettext('Triggered') : $gettext('OK') }}
          </Tag>
        </div>
      </Space>
    </Card>
  </div>
</template>
