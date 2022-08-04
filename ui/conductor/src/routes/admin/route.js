export default {
  name: 'admin',
  path: '/admin',
  component: () => import('./Admin.vue'),
  meta: {
    title: 'Administrator'
  }
}
