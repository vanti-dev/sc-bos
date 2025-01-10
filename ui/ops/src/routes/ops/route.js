import SidebarPage from '@/components/pages/SidebarPage.vue';
import notifications from '@/routes/ops/notifications/route.js';
import overview from '@/routes/ops/overview/route.js';
import {useUiConfigStore} from '@/stores/uiConfig.js';

import {route} from '@/util/router.js';

export default {
  name: 'ops',
  path: '/ops',
  redirect: () => {
    const uiConfig = useUiConfigStore();
    if (uiConfig.pathEnabled('/ops/overview')) {
      return '/ops/overview';
    } else if (uiConfig.pathEnabled('/ops/notifications')) {
      return '/ops/notifications';
    } else if (uiConfig.pathEnabled('/ops/air-quality')) {
      return '/ops/air-quality';
    } else if (uiConfig.pathEnabled('/ops/emergency-lighting')) {
      return '/ops/emergency-lighting';
    } else if (uiConfig.pathEnabled('/ops/security')) {
      return '/ops/security';
    }
    return '/ops/loading';
  },
  components: {
    default: SidebarPage,
    nav: () => import('./OpsNav.vue')
  },
  children: [
    {
      path: 'loading',
      component: () => import('./OpsLoading.vue')
    },
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
    const appConfig = useUiConfigStore();
    await appConfig.loadConfig();
    if (to.path === '/ops') {
      if (appConfig.pathEnabled('/ops/overview')) {
        next('/ops/overview');
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
