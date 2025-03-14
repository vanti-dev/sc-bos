export default [
  {
    path: 'roles/:roleId?',
    name: 'roles',
    components: {
      default: () => import('./RolesPage.vue'),
      sidebar: () => import('./RolesSideBar.vue')
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
