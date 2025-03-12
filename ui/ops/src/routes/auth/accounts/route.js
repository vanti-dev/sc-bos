export default [
  {
    path: 'accounts/:accountId?',
    components: {
      default: () => import('./AccountsPage.vue')
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
