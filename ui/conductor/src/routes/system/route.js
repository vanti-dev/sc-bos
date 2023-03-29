import SidebarPage from '@/components/page-layout/SidebarPage.vue';
import {featureEnabled} from '@/routes/config';

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
    }
  ],
  meta: {
    title: 'System'
  },
  beforeEnter: async (to, from, next) => {
    if (to.path === '/system') {
      if (await featureEnabled('/system/drivers')) {
        next('/system/drivers');
      } else if (await featureEnabled('/system/features')) {
        next('/system/features');
      }
    } else {
      next();
    }
  }
};
