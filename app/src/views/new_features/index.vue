<script setup lang="ts">
import type { EntryItem } from '@/components/EntryGrid/types'
import {
  AlertOutlined,
  AreaChartOutlined,
  FileImageOutlined,
  HeartOutlined,
  SafetyCertificateOutlined,
  ThunderboltOutlined,
} from '@antdv-next/icons'
import EntryGrid from '@/components/EntryGrid/EntryGrid.vue'

const entries = computed<EntryItem[]>(() => [
  {
    title: $gettext('Domain Analytics'),
    description: $gettext('Traffic, bandwidth, bots vs humans, and security bounce stats.'),
    path: '/new_features/domain_analytics',
    icon: AreaChartOutlined,
  },
  {
    title: $gettext('Media Manager'),
    description: $gettext('Browse, upload, and manage image folders served by Nginx.'),
    path: '/new_features/media',
    icon: FileImageOutlined,
  },
  {
    title: $gettext('Site Shield'),
    description: $gettext('Hotlink protection, rate limits, WAF lite, IP blocks, cache headers.'),
    path: '/new_features/site_shield',
    icon: SafetyCertificateOutlined,
  },
  {
    title: $gettext('Image Site Kit'),
    description: $gettext('One-click comic/card image hosting security + cache preset.'),
    path: '/new_features/image_site_kit',
    icon: ThunderboltOutlined,
  },
  {
    title: $gettext('Origin Health'),
    description: $gettext('Uptime, latency, and failure status for tracked site origins.'),
    path: '/new_features/origin_health',
    icon: HeartOutlined,
  },
  {
    title: $gettext('Traffic Alerts'),
    description: $gettext('Notify when requests, bandwidth, or blocks spike past a threshold.'),
    path: '/new_features/traffic_alerts',
    icon: AlertOutlined,
  },
])
</script>

<template>
  <EntryGrid :entries="entries" />
</template>
