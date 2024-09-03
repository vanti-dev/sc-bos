export default [
  {
    path: 'third-party/:accountId?',
    components: {
      default: () => import('./AccountList.vue'),
      sidebar: () => import('./AccountSideBar.vue')
    },
    props: {
      default: false,
      sidebar: true
    },
    meta: {
      authentication: {
        rolesRequired: ['superAdmin', 'admin', 'operator', 'viewer']
      }
    }
  }
];
