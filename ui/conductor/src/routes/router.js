import auth from '@/routes/auth/route.js';
import devices from '@/routes/devices/route.js';
import commission from '@/routes/commission/route.js';
import design from '@/routes/design/route.js';
import ops from '@/routes/ops/route.js';
import start from '@/routes/start/route.js';
import {route, routeTitle} from '@/util/router.js';
import Vue, {nextTick} from 'vue';
import VueRouter from 'vue-router';

Vue.use(VueRouter);

const router = new VueRouter({
  mode: 'history',
  routes: [
    ...route(auth),
    ...route(devices),
    ...route(ops),
    ...route(design),
    ...route(commission),
  ]
});

if (window) {
  router.afterEach((to, from) => {
    const nt = routeTitle(to);
    const ot = routeTitle(from);
    if (nt === ot) {
      return;
    }

    const title = nt ? `${nt} - Smart Core` : `Smart Core`;
    nextTick(() => window.document.title = title);
  })
}

export default router;
