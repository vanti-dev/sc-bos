import ChildOnlyPage from '../../components/ChildOnlyPage.vue';

export default {
  name: 'commission',
  path: '/commission',
  redirect: '/commission/spaces',
  components: {
    default: ChildOnlyPage,
    sections: () => import('./CommissionSections.vue'),
  },
  children: [
    {path: 'spaces', component: () => import('./CommissionSpaces.vue')},
    {path: 'devices', component: () => import('./CommissionDevices.vue')},
  ],
  meta: {
    title: 'Commissioner'
  }
}
