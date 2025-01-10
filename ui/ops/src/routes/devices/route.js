import SidebarPage from '@/components/pages/SidebarPage.vue';

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
      },
      meta: {
        authentication: {
          rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
        }
      }
    }
  ],
  meta: {
    title: 'Devices',
    authentication: {
      rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
    }
  }
};
