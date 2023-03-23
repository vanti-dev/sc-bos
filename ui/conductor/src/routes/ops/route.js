import ChildOnlyPage from '@/components/page-layout/ChildOnlyPage.vue';
import notifications from '@/routes/ops/notifications/route.js';

import {route} from '@/util/router.js';

export default {
  name: 'ops',
  path: '/ops',
  redirect: '/ops/overview',
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
  }
};
