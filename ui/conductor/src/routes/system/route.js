import SidebarPage from '@/layout/SidebarPage.vue';
import {useAppConfigStore} from '@/stores/app-config';

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
      }
    },
    {
      path: 'features',
      components: {
        default: () => import('./components/pages/FeaturesList.vue'),
        sidebar: () => import('./components/ServicesSideBar.vue')
      }
    },
    {
      path: 'components',
      components: {
        default: () => import('./components/pages/ComponentsList.vue')
      }
    }
  ],
  meta: {
    title: 'System'
  },
  beforeEnter: async (to, from, next) => {
    const appConfig = useAppConfigStore();
    if (to.path === '/system') {
      if (appConfig.pathEnabled('/system/drivers')) {
        next('/system/drivers');
      } else if (appConfig.pathEnabled('/system/features')) {
        next('/system/features');
      }
    } else {
      next();
    }
  }
};
