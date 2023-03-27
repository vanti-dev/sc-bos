import SidebarPage from '@/components/page-layout/SidebarPage.vue';

export default {
  name: 'system',
  path: '/system',
  components: {
    default: SidebarPage,
    nav: () => import('./SystemNav.vue')
  },
  children: [
    {
      path: 'drivers',
      components: {
        default: () => import('./components/ServicesList.vue')
      }
    },
    {
      path: 'features',
      components: {
        default: () => import('./components/ServicesList.vue')
      }
    }
  ],
  meta: {
    title: 'System'
  }
};
