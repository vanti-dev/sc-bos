import admin from '@/routes/admin/route.js';
import commission from '@/routes/commission/route.js';
import design from '@/routes/design/route.js';
import operate from '@/routes/operate/route.js';
import start from '@/routes/start/route.js';
import {route, routeTitle} from '@/util/router.js';
import Vue, {nextTick} from 'vue';
import VueRouter from 'vue-router';

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
