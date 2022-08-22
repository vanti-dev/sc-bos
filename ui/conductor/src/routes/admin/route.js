import ChildOnlyPage from '../../components/ChildOnlyPage.vue';

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
    {path: 'tenants', component: () => import('./AdminTenants.vue')}
  ],
  meta: {
    title: 'Administrator'
  }
}
