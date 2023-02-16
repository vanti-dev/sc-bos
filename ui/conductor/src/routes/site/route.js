import ChildOnlyPage from '@/components/page-layout/ChildOnlyPage.vue';

export default {
  name: 'site',
  path: '/site',
  redirect: '/site/zone',
  components: {
    default: ChildOnlyPage,
    nav: () => import('./ZoneNav.vue')
  },
  children: [
    {
      name: 'zone',
      path: 'zone',
      children: [
        {path: 'list', component: () => import('./ZoneList.vue')}
      ]
    }
  ],
  meta: {
    title: 'Site Settings'
  }
};
