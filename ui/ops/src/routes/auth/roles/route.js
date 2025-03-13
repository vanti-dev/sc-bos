export default [
  {
    path: 'roles/:roleId?',
    components: {
      default: () => import('./RolesPage.vue')
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
