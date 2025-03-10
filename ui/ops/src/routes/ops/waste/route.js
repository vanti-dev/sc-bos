export default [
  {
    path: 'waste',
    components: {
      default: () => import('./WasteTable.vue')
    },
    props: {
      default: true,
      sidebar: false
    },
    meta: {
      authentication: {
        rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
      }
    }
  }
];