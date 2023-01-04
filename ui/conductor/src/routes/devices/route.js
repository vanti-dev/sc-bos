import SidebarPage from '@/components/page-layout/SidebarPage.vue';
import {route} from '@/util/router.js';
import lighting from '@/routes/devices/lighting/route';

export default {
  name: 'devices',
  path: '/devices',
  redirect: '/devices/lighting',
  components: {
    default: SidebarPage,
    nav: () => import('./DevicesNav.vue')
  },
  children: [
    ...route(lighting)
  ],
  meta: {
    title: 'Devices'
  }
};
