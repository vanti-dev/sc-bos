import ChildOnlyPage from '@/components/ChildOnlyPage.vue';
import tenants from '@/routes/admin/tenant/route.js';
import {route} from '@/util/router.js';

export default {
  name: 'admin',
  path: '/admin',
  redirect: '/admin/summary',
  components: {
    default: ChildOnlyPage,
    nav: () => import('./AdminNav.vue')
  },
  children: [
    {path: 'summary', component: () => import('./AdminSummary.vue')},
      ...route(tenants)
  ],
  meta: {
    title: 'Administrator'
  }
}
