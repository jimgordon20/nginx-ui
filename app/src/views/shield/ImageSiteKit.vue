<script setup lang="ts">
import type { Site } from '@/api/site'
import { ThunderboltOutlined } from '@antdv-next/icons'
import { Alert, App, Button, Card, Select, Space, Typography } from 'antdv-next'
import siteApi from '@/api/site'
import shieldApi from '@/api/shield'

const { message, modal } = App.useApp()

const sites = ref<Site[]>([])
const selectedSite = ref('')
const applying = ref(false)

async function loadSites() {
  const res = await siteApi.getList({ page: 1, page_size: 500 }) as { data?: Site[] }
  sites.value = Array.isArray(res?.data) ? res.data : []
  if (!selectedSite.value && sites.value.length)
    selectedSite.value = sites.value[0].name
}

function applyPreset() {
  if (!selectedSite.value)
    return
  modal.confirm({
    title: $gettext('Apply Image Site Kit'),
    content: $gettext('Enables hotlink protection, image rate limits, WAF lite, security headers, and long-lived cache for static media on this site.'),
    okText: $gettext('Apply'),
    async onOk() {
      applying.value = true
      try {
        await shieldApi.applyImagePreset(selectedSite.value)
        message.success($gettext('Image Site Kit applied'))
      }
      catch (error) {
        console.error(error)
        message.error($gettext('Failed to apply Image Site Kit'))
        throw error
      }
      finally {
        applying.value = false
      }
    },
  })
}

onMounted(() => {
  void loadSites()
})
</script>

<template>
  <Card :title="$gettext('Image Site Kit')" variant="borderless">
    <Alert
      class="mb-4"
      type="info"
      show-icon
      :message="$gettext('One-click comic / trading-card image hosting preset')"
      :description="$gettext('Best starting point for image-heavy sites. You can fine-tune everything afterward in Site Shield.')"
    />

    <Typography.Paragraph>
      {{ $gettext('This preset turns on:') }}
    </Typography.Paragraph>
    <ul class="mb-4 opacity-80">
      <li>{{ $gettext('Hotlink protection for images') }}</li>
      <li>{{ $gettext('Per-image rate limiting') }}</li>
      <li>{{ $gettext('WAF lite (empty UA + scrapers + sensitive paths)') }}</li>
      <li>{{ $gettext('Security headers') }}</li>
      <li>{{ $gettext('Long cache headers for static media') }}</li>
    </ul>

    <Space wrap>
      <Select
        v-model:value="selectedSite"
        class="min-w-56"
        show-search
        :options="sites.map(site => ({ label: site.name, value: site.name }))"
        :placeholder="$gettext('Select site')"
      />
      <Button
        type="primary"
        :loading="applying"
        :disabled="!selectedSite"
        @click="applyPreset"
      >
        <template #icon>
          <ThunderboltOutlined />
        </template>
        {{ $gettext('Apply Image Site Kit') }}
      </Button>
    </Space>
  </Card>
</template>

<style scoped>
.min-w-56 {
  min-width: 14rem;
}
</style>
