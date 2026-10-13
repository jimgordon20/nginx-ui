import type { RouteRecordRaw } from 'vue-router'
import { AppstoreOutlined, AreaChartOutlined, FileImageOutlined } from '@antdv-next/icons'

export const newFeaturesRoutes: RouteRecordRaw[] = [
  {
    path: 'new_features',
    name: 'New Features',
    component: () => import('@/layouts/BaseRouterView.vue'),
    meta: {
      name: () => $gettext('New Features'),
      icon: AppstoreOutlined,
    },
    children: [
      {
        path: '',
        name: 'New Features Home',
        component: () => import('@/views/new_features/index.vue'),
        meta: {
          name: () => $gettext('New Features'),
          hiddenInSidebar: true,
        },
      },
      {
        path: 'domain_analytics',
        name: 'Domain Analytics',
        component: () => import('@/views/domain_analytics/DomainAnalytics.vue'),
        meta: {
          name: () => $gettext('Domain Analytics'),
          icon: AreaChartOutlined,
        },
      },
      {
        path: 'media',
        name: 'Media Manager',
        component: () => import('@/views/media/MediaManager.vue'),
        meta: {
          name: () => $gettext('Media Manager'),
          icon: FileImageOutlined,
        },
      },
    ],
  },
]
