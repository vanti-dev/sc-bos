export default [
  {
    path: 'security-events',
    components: {
      default: () => import('./SecurityEventsTable.vue')

    },
    props: {
      default: true,
      sidebar: false
    },
    meta: {
      authentication: {
        rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator']
      }
    }
  }
];
