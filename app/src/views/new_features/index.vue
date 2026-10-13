<script setup lang="ts">
import type { EntryItem } from '@/components/EntryGrid/types'
import { AreaChartOutlined, FileImageOutlined } from '@antdv-next/icons'
import EntryGrid from '@/components/EntryGrid/EntryGrid.vue'

const entries = computed<EntryItem[]>(() => [
  {
    title: $gettext('Domain Analytics'),
    description: $gettext('Cloudflare-style traffic analytics for domain access logs.'),
    path: '/new_features/domain_analytics',
    icon: AreaChartOutlined,
  },
  {
    title: $gettext('Media Manager'),
    description: $gettext('Browse, upload, and manage image folders served by Nginx.'),
    path: '/new_features/media',
    icon: FileImageOutlined,
  },
])
</script>

<template>
  <EntryGrid :entries="entries" />
</template>
