import type { RouteRecordRaw } from 'vue-router'
import { AreaChartOutlined } from '@antdv-next/icons'

export const domainAnalyticsRoutes: RouteRecordRaw[] = [
  {
    path: 'domain_analytics',
    name: 'Domain Analytics',
    component: () => import('@/views/domain_analytics/DomainAnalytics.vue'),
    meta: {
      name: () => $gettext('Domain Analytics'),
      icon: AreaChartOutlined,
    },
  },
]
