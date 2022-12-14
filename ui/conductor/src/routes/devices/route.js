import ChildOnlyPage from '@/components/ChildOnlyPage.vue';
import {route} from '@/util/router.js';
import lighting from '@/routes/devices/lighting/route'

export default {
  name: 'devices',
  path: '/devices',
  redirect: '/devices/lighting',
  components: {
    default: ChildOnlyPage,
    nav: () => import('./DevicesNav.vue')
  },
  children: [
    ...route(lighting)
  ],
  meta: {
    title: 'Devices'
  }
}
