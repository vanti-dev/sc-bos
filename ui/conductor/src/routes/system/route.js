import SidebarPage from '@/components/page-layout/SidebarPage.vue';
import {useUiConfigStore} from '@/stores/uiConfig.js';

export default {
  name: 'system',
  path: '/system',
  components: {
    default: SidebarPage,
    nav: () => import('./SystemNav.vue')
  },
  children: [
    {
      path: 'drivers',
      components: {
        default: () => import('./components/pages/DriversList.vue'),
        sidebar: () => import('./components/ServicesSideBar.vue')
      },
      meta: {
        authentication: {
          rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
        }
      }
    },
    {
      path: 'features',
      components: {
        default: () => import('./components/pages/FeaturesList.vue'),
        sidebar: () => import('./components/ServicesSideBar.vue')
      },
      meta: {
        authentication: {
          rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
        }
      }
    },
    {
      path: 'components',
      components: {
        default: () => import('./components/pages/ComponentsList.vue')
      },
      meta: {
        authentication: {
          rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
        }
      }
    }
  ],
  meta: {
    authentication: {
      rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
    },
    title: 'System'
  },
  beforeEnter: async (to, from, next) => {
    const uiConfig = useUiConfigStore();
    if (to.path === '/system') {
      if (uiConfig.pathEnabled('/system/drivers')) {
        next('/system/drivers');
      } else if (uiConfig.pathEnabled('/system/features')) {
        next('/system/features');
      }
    } else {
      next();
    }
  }
};
