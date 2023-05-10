import ChildOnlyPage from '@/components/page-layout/ChildOnlyPage.vue';
import notifications from '@/routes/ops/notifications/route.js';
import {useAppConfigStore} from '@/stores/app-config';

import {route} from '@/util/router.js';

export default {
  name: 'ops',
  path: '/ops',
  components: {
    default: ChildOnlyPage,
    nav: () => import('./OpsNav.vue')
  },
  children: [
    {path: 'overview', component: () => import('./OpsHome.vue')},
    {path: 'emergency-lighting', component: () => import('./emergency-lighting/EmergencyLighting.vue')},
    {path: 'security', component: () => import('./security/SecurityHome.vue')},
    ...route(notifications)
  ],
  meta: {
    title: 'Operations'
  },
  beforeEnter: async (to, from, next) => {
    const appConfig = useAppConfigStore();
    if (to.path === '/ops') {
      if (appConfig.pathEnabled('/ops/overview')) {
        next('/ops/overview');
      } else if (appConfig.pathEnabled('/ops/notifications')) {
        next('/ops/notifications');
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
