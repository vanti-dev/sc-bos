export default [
  {
    path: 'hvac',
    components: {
      default: () => import('./HvacTable.vue'),
      sidebar: () => import('./HvacSideBar.vue')
    }
  }
];
