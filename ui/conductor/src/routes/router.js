import VueRouter from 'vue-router';
import admin from './admin/route.js';
import commission from './commission/route.js';
import design from './design/route.js';
import operate from './operate/route.js';
import start from './start/route.js';
import Vue, {nextTick} from 'vue';
import {route, routeTitle} from '../util/router.js';

Vue.use(VueRouter);

const router = new VueRouter({
  mode: 'history',
  routes: [
    ...route(start),
    ...route(design),
    ...route(commission),
    ...route(operate),
    ...route(admin),
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
