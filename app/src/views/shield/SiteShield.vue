<script setup lang="ts">
import type { Site } from '@/api/site'
import type { ShieldConfig } from '@/api/shield'
import { ReloadOutlined, SafetyCertificateOutlined } from '@antdv-next/icons'
import {
  App,
  Button,
  Card,
  Col,
  Divider,
  Input,
  InputNumber,
  Row,
  Select,
  Space,
  Switch,
  Tag,
} from 'antdv-next'
import siteApi from '@/api/site'
import shieldApi from '@/api/shield'

const { message, modal } = App.useApp()

const loading = ref(false)
const saving = ref(false)
const sites = ref<Site[]>([])
const selectedSite = ref('')
const config = ref<ShieldConfig | null>(null)
const allowReferersText = ref('')
const denyIpsText = ref('')
const blockedAgentsText = ref('')
const blockedPathsText = ref('')

async function loadSites() {
  try {
    const res = await siteApi.getList({ page: 1, page_size: 500 }) as { data?: Site[] }
    sites.value = Array.isArray(res?.data) ? res.data : []
    if (!selectedSite.value && sites.value.length)
      selectedSite.value = sites.value[0].name
  }
  catch (error) {
    console.error(error)
    message.error($gettext('Failed to load sites'))
  }
}

async function loadConfig() {
  if (!selectedSite.value)
    return
  loading.value = true
  try {
    const cfg = await shieldApi.get(selectedSite.value)
    config.value = cfg
    allowReferersText.value = (cfg.hotlink.allow_referers || []).join(', ')
    denyIpsText.value = (cfg.ip_block.deny_ips || []).join('\n')
    blockedAgentsText.value = (cfg.waf.blocked_user_agents || []).join(', ')
    blockedPathsText.value = (cfg.waf.blocked_paths || []).join('\n')
  }
  catch (error) {
    console.error(error)
    message.error($gettext('Failed to load Site Shield settings'))
  }
  finally {
    loading.value = false
  }
}

function syncTextFields() {
  if (!config.value)
    return
  config.value.hotlink.allow_referers = allowReferersText.value
    .split(',')
    .map(item => item.trim())
    .filter(Boolean)
  config.value.ip_block.deny_ips = denyIpsText.value
    .split(/[\n,]+/)
    .map(item => item.trim())
    .filter(Boolean)
  config.value.waf.blocked_user_agents = blockedAgentsText.value
    .split(',')
    .map(item => item.trim())
    .filter(Boolean)
  config.value.waf.blocked_paths = blockedPathsText.value
    .split(/[\n,]+/)
    .map(item => item.trim())
    .filter(Boolean)
}

async function saveOnly() {
  if (!config.value || !selectedSite.value)
    return
  syncTextFields()
  saving.value = true
  try {
    config.value = await shieldApi.save(selectedSite.value, config.value)
    message.success($gettext('Saved'))
  }
  catch (error) {
    console.error(error)
    message.error($gettext('Save failed'))
  }
  finally {
    saving.value = false
  }
}

async function applyNow() {
  if (!config.value || !selectedSite.value)
    return
  syncTextFields()
  modal.confirm({
    title: $gettext('Apply Site Shield'),
    content: $gettext('This writes managed nginx rules into the selected site and reloads Nginx.'),
    okText: $gettext('Apply'),
    async onOk() {
      saving.value = true
      try {
        config.value = await shieldApi.apply(selectedSite.value, config.value!)
        message.success($gettext('Site Shield applied'))
      }
      catch (error) {
        console.error(error)
        message.error($gettext('Apply failed — check nginx test output'))
        throw error
      }
      finally {
        saving.value = false
      }
    },
  })
}

async function detach() {
  if (!selectedSite.value)
    return
  modal.confirm({
    title: $gettext('Remove Site Shield'),
    content: $gettext('Remove managed shield includes from this site config?'),
    okType: 'danger',
    async onOk() {
      config.value = await shieldApi.detach(selectedSite.value)
      message.success($gettext('Detached'))
    },
  })
}

watch(selectedSite, () => {
  void loadConfig()
})

onMounted(async () => {
  await loadSites()
  await loadConfig()
})
</script>

<template>
  <div>
    <Card
      :title="$gettext('Site Shield')"
      variant="borderless"
      class="mb-4"
    >
      <template #extra>
        <Space wrap>
          <Select
            v-model:value="selectedSite"
            class="min-w-56"
            :options="sites.map(site => ({ label: site.name, value: site.name }))"
            :placeholder="$gettext('Select site')"
            show-search
          />
          <Button :loading="loading" @click="loadConfig">
            <template #icon>
              <ReloadOutlined />
            </template>
          </Button>
          <Button :loading="saving" @click="saveOnly">
            {{ $gettext('Save') }}
          </Button>
          <Button type="primary" :loading="saving" @click="applyNow">
            <template #icon>
              <SafetyCertificateOutlined />
            </template>
            {{ $gettext('Apply to Site') }}
          </Button>
          <Button danger :disabled="!config?.applied" @click="detach">
            {{ $gettext('Detach') }}
          </Button>
        </Space>
      </template>

      <div class="mb-3">
        <Tag v-if="config?.applied" color="success">
          {{ $gettext('Applied') }}
        </Tag>
        <Tag v-else>
          {{ $gettext('Not applied') }}
        </Tag>
        <span class="opacity-70 text-sm ml-2">
          {{ $gettext('Cloudflare-style protection for hotlinks, bots, cache, and abuse.') }}
        </span>
      </div>
    </Card>

    <div v-if="config">
      <Row :gutter="[16, 16]">
        <Col :xs="24" :lg="12">
          <Card :title="$gettext('Hotlink Protection')" variant="borderless">
            <div class="flex items-center justify-between mb-3">
              <span>{{ $gettext('Enabled') }}</span>
              <Switch v-model:checked="config.hotlink.enabled" />
            </div>
            <div class="flex items-center justify-between mb-3">
              <span>{{ $gettext('Allow empty referer') }}</span>
              <Switch v-model:checked="config.hotlink.allow_empty" :disabled="!config.hotlink.enabled" />
            </div>
            <div class="mb-1 text-sm opacity-70">
              {{ $gettext('Extra allowed referers (comma separated)') }}
            </div>
            <Input
              v-model:value="allowReferersText"
              :disabled="!config.hotlink.enabled"
              :placeholder="$gettext('*.mysite.com example.com')"
            />
          </Card>
        </Col>

        <Col :xs="24" :lg="12">
          <Card :title="$gettext('Rate Limit')" variant="borderless">
            <div class="flex items-center justify-between mb-3">
              <span>{{ $gettext('Enabled') }}</span>
              <Switch v-model:checked="config.rate_limit.enabled" />
            </div>
            <div class="mb-1 text-sm opacity-70">
              {{ $gettext('Rate') }}
            </div>
            <Input v-model:value="config.rate_limit.rate" class="mb-3" :disabled="!config.rate_limit.enabled" />
            <div class="mb-1 text-sm opacity-70">
              {{ $gettext('Burst') }}
            </div>
            <InputNumber v-model:value="config.rate_limit.burst" class="w-full" :disabled="!config.rate_limit.enabled" />
          </Card>
        </Col>

        <Col :xs="24" :lg="12">
          <Card :title="$gettext('WAF Lite')" variant="borderless">
            <div class="flex items-center justify-between mb-3">
              <span>{{ $gettext('Enabled') }}</span>
              <Switch v-model:checked="config.waf.enabled" />
            </div>
            <div class="flex items-center justify-between mb-3">
              <span>{{ $gettext('Block empty User-Agent') }}</span>
              <Switch v-model:checked="config.waf.block_empty_ua" :disabled="!config.waf.enabled" />
            </div>
            <div class="flex items-center justify-between mb-3">
              <span>{{ $gettext('Block bad bots') }}</span>
              <Switch v-model:checked="config.waf.block_bad_bots" :disabled="!config.waf.enabled" />
            </div>
            <div class="mb-1 text-sm opacity-70">
              {{ $gettext('Blocked agents') }}
            </div>
            <Input
              v-model:value="blockedAgentsText"
              class="mb-3"
              :disabled="!config.waf.enabled"
            />
            <div class="mb-1 text-sm opacity-70">
              {{ $gettext('Blocked paths (one per line)') }}
            </div>
            <ATextarea
              v-model:value="blockedPathsText"
              :rows="3"
              :disabled="!config.waf.enabled"
            />
          </Card>
        </Col>

        <Col :xs="24" :lg="12">
          <Card :title="$gettext('Under Attack Mode')" variant="borderless">
            <div class="flex items-center justify-between mb-3">
              <span>{{ $gettext('Enabled') }}</span>
              <Switch v-model:checked="config.under_attack.enabled" />
            </div>
            <div class="mb-1 text-sm opacity-70">
              {{ $gettext('Strict rate') }}
            </div>
            <Input
              v-model:value="config.under_attack.rate"
              class="mb-3"
              :disabled="!config.under_attack.enabled"
            />
            <div class="mb-1 text-sm opacity-70">
              {{ $gettext('Burst') }}
            </div>
            <InputNumber
              v-model:value="config.under_attack.burst"
              class="w-full"
              :disabled="!config.under_attack.enabled"
            />
          </Card>
        </Col>

        <Col :xs="24" :lg="12">
          <Card :title="$gettext('IP Blocklist')" variant="borderless">
            <div class="flex items-center justify-between mb-3">
              <span>{{ $gettext('Enabled') }}</span>
              <Switch v-model:checked="config.ip_block.enabled" />
            </div>
            <div class="mb-1 text-sm opacity-70">
              {{ $gettext('Deny IPs (one per line)') }}
            </div>
            <ATextarea
              v-model:value="denyIpsText"
              :rows="4"
              :disabled="!config.ip_block.enabled"
              placeholder="1.2.3.4&#10;5.6.7.0/24"
            />
          </Card>
        </Col>

        <Col :xs="24" :lg="12">
          <Card :title="$gettext('Cache & CDN Headers')" variant="borderless">
            <div class="flex items-center justify-between mb-3">
              <span>{{ $gettext('Security headers') }}</span>
              <Switch v-model:checked="config.security_headers.enabled" />
            </div>
            <Divider class="my-3" />
            <div class="flex items-center justify-between mb-3">
              <span>{{ $gettext('Image/static cache') }}</span>
              <Switch v-model:checked="config.cache.enabled" />
            </div>
            <div class="mb-1 text-sm opacity-70">
              {{ $gettext('Expires') }}
            </div>
            <Input v-model:value="config.cache.expires" class="mb-3" :disabled="!config.cache.enabled" />
            <div class="mb-1 text-sm opacity-70">
              {{ $gettext('Cache-Control') }}
            </div>
            <Input v-model:value="config.cache.cache_control" :disabled="!config.cache.enabled" />
          </Card>
        </Col>
      </Row>

      <Card
        v-if="config.preview"
        class="mt-4"
        :title="$gettext('Generated Nginx Preview')"
        variant="borderless"
      >
        <pre class="preview-code">{{ config.preview }}</pre>
      </Card>
    </div>
  </div>
</template>

<style scoped>
.preview-code {
  margin: 0;
  padding: 12px;
  border-radius: 8px;
  background: var(--ant-color-fill-tertiary);
  color: var(--ant-color-text);
  font-size: 12px;
  overflow: auto;
  max-height: 360px;
}

.min-w-56 {
  min-width: 14rem;
}
</style>
