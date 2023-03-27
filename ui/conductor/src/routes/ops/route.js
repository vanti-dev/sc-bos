import ChildOnlyPage from '@/components/page-layout/ChildOnlyPage.vue';
import notifications from '@/routes/ops/notifications/route.js';

import {route} from '@/util/router.js';
import {featureEnabled} from '@/routes/config';

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
    ...route(notifications)
  ],
  meta: {
    title: 'Operations'
  },
  beforeEnter: async (to, from, next) => {
    if (to.path === '/ops') {
      if (await featureEnabled('/ops/overview')) {
        next('/ops/overview');
      } else if (await featureEnabled('/ops/notifications')) {
        next('/ops/notifications');
      } else if (await featureEnabled('/ops/emergency-lighting')) {
        next('/ops/emergency-lighting');
      }
    } else {
      next();
    }
  }
};
