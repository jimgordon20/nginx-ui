import type { RouteRecordRaw } from 'vue-router'
import {
  AlertOutlined,
  AppstoreOutlined,
  AreaChartOutlined,
  FileImageOutlined,
  HeartOutlined,
  SafetyCertificateOutlined,
  ThunderboltOutlined,
} from '@antdv-next/icons'

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
      {
        path: 'site_shield',
        name: 'Site Shield',
        component: () => import('@/views/shield/SiteShield.vue'),
        meta: {
          name: () => $gettext('Site Shield'),
          icon: SafetyCertificateOutlined,
        },
      },
      {
        path: 'image_site_kit',
        name: 'Image Site Kit',
        component: () => import('@/views/shield/ImageSiteKit.vue'),
        meta: {
          name: () => $gettext('Image Site Kit'),
          icon: ThunderboltOutlined,
        },
      },
      {
        path: 'origin_health',
        name: 'Origin Health',
        component: () => import('@/views/shield/OriginHealth.vue'),
        meta: {
          name: () => $gettext('Origin Health'),
          icon: HeartOutlined,
        },
      },
      {
        path: 'traffic_alerts',
        name: 'Traffic Alerts',
        component: () => import('@/views/shield/TrafficAlerts.vue'),
        meta: {
          name: () => $gettext('Traffic Alerts'),
          icon: AlertOutlined,
        },
      },
    ],
  },
]
