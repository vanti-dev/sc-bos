export default [
  {
    path: 'notifications',
    components: {
      default: () => import('./Notifications.vue'),
      sidebar: () => import('./NotificationSideBar.vue')
    },
    props: {
      default: true,
      sidebar: false
    }
  }
];
