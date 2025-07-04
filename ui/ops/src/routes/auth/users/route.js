export default [
  {
    path: 'users/:accountId?',
    name: 'users',
    components: {
      default: () => import('./UsersPage.vue'),
      sidebar: () => import('./UsersSideBar.vue')
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
