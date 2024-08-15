import ChildOnlyPage from '@/components/page-layout/ChildOnlyPage.vue';

export default {
  name: 'site',
  path: '/site',
  redirect: '/site/zone',
  components: {
    default: ChildOnlyPage,
    nav: () => import('./SiteNav.vue')
  },
  children: [
    {
      name: 'zone',
      path: 'zone/:zone(.*)*',
      component: () => import('./zone/ZonePage.vue'),
      props: true,
      meta: {
        authentication: {
          rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
        }
      }
    }
  ],
  meta: {
    authentication: {
      rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
    },
    title: 'Site Config'
  }
};
