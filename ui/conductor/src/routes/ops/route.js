import ChildOnlyPage from '@/components/page-layout/ChildOnlyPage.vue';
import notifications from '@/routes/ops/notifications/route.js';

import {route} from '@/util/router.js';

export default {
  name: 'ops',
  path: '/ops',
  redirect: '/ops/summary',
  components: {
    default: ChildOnlyPage,
    nav: () => import('./OpsNav.vue')
  },
  children: [
    {path: 'summary', component: () => import('./OpsSummary.vue')},
    ...route(notifications)
  ],
  meta: {
    title: 'Operations'
  }
};
