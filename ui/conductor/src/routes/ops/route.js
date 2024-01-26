import SidebarPage from '@/components/page-layout/SidebarPage.vue';
import overview from '@/routes/ops/overview/route.js';
import notifications from '@/routes/ops/notifications/route.js';
import {useAppConfigStore} from '@/stores/app-config';

import {route} from '@/util/router.js';

export default {
  name: 'ops',
  path: '/ops',
  components: {
    default: SidebarPage,
    nav: () => import('./OpsNav.vue')
  },
  children: [
    ...route(overview),
    {
      path: 'emergency-lighting',
      component: () => import('./emergency-lighting/EmergencyLighting.vue'),
      meta: {
        authentication: {
          rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
        }
      }
    },
    {
      path: 'security',
      component: () => import('./security/SecurityHome.vue'),
      meta: {
        authentication: {
          rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
        }
      }
    },
    {
      path: 'air-quality',
      component: () => import('./air-quality/AirQuality.vue'),
      meta: {
        authentication: {
          rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
        }
      }
    },
    ...route(notifications)
  ],
  meta: {
    authentication: {
      rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
    },
    title: 'Operations'
  },
  beforeEnter: async (to, from, next) => {
    const appConfig = useAppConfigStore();
    if (to.path === '/ops') {
      if (appConfig.pathEnabled('/ops/overview')) {
        next('/ops/overview/building');
      } else if (appConfig.pathEnabled('/ops/notifications')) {
        next('/ops/notifications');
      } else if (appConfig.pathEnabled('/ops/air-quality')) {
        next('/ops/air-quality');
      } else if (appConfig.pathEnabled('/ops/emergency-lighting')) {
        next('/ops/emergency-lighting');
      } else if (appConfig.pathEnabled('/ops/security')) {
        next('/ops/security');
      }
    } else {
      next();
    }
  }
};
