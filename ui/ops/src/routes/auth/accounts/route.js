export default [
  {
    path: 'accounts/:accountId?',
    name: 'accounts',
    components: {
      default: () => import('./AccountsPage.vue'),
      sidebar: () => import('./AccountsSideBar.vue')
    },
    props: {
      default: true,
      sidebar: false
    },
    meta: {
      authentication: {
        rolesRequired: ['superAdmin', 'admin', 'operator', 'viewer']
      }
    }
  }
];
