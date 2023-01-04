export default [
  {
    path: 'lighting',
    components: {
      default: () => import('./LightingTable.vue'),
      sidebar: () => import('./RowMenu.vue')
    }
  }
];
