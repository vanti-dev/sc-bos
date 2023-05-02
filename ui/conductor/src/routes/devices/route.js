import SidebarPage from '@/layout/SidebarPage.vue';

export default {
  name: 'devices',
  path: '/devices',
  redirect: '/devices/all',
  components: {
    default: SidebarPage,
    nav: () => import('./DevicesNav.vue')
  },
  children: [
    {
      path: ':subsystem',
      components: {
        default: () => import('./components/DeviceTable.vue'),
        sidebar: () => import('./components/DeviceSideBar.vue')
      },
      props: {
        default: true,
        sidebar: false
      }
    }
  ],
  meta: {
    title: 'Devices'
  }
};
