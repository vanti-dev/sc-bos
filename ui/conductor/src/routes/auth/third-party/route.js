import SidebarPage from '@/components/page-layout/SidebarPage.vue';

export default [
  {
    name: 'third-party',
    path: '/auth/third-party',
    components: {
      default: SidebarPage,
      nav: () => import('../AuthNav.vue')
    },
    children: [
      {
        path: ':accountId?',
        components: {
          default: () => import('./AccountList.vue'),
          sidebar: () => import('./AccountSideBar.vue')
        },
        props: {
          default: false,
          sidebar: true
        }
      }
    ],
    meta: {
      title: 'Auth'
    }
  }
];
