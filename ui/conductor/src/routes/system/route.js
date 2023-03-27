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
        default: () => import('./components/pages/DriversList.vue'),
        sidebar: () => import('./components/ServicesSideBar.vue')
      }
    },
    {
      path: 'features',
      components: {
        default: () => import('./components/pages/FeaturesList.vue'),
        sidebar: () => import('./components/ServicesSideBar.vue')
      }
    }
  ],
  meta: {
    title: 'System'
  }
};
