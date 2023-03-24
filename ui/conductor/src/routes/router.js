import auth from '@/routes/auth/route.js';
import devices from '@/routes/devices/route.js';
import ops from '@/routes/ops/route.js';
import automations from '@/routes/automations/route.js';
import site from '@/routes/site/route.js';
import {route, routeTitle} from '@/util/router.js';
import Vue, {nextTick} from 'vue';
import VueRouter from 'vue-router';
import {featureEnabled} from '@/routes/config';

Vue.use(VueRouter);

const router = new VueRouter({
  mode: 'history',
  routes: [
    {path: '/', redirect: '/ops'},
    ...route(auth),
    ...route(devices),
    ...route(ops),
    ...route(automations),
    ...route(site)
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
  });
  router.beforeEach(async (to, from, next) => {
    next(await featureEnabled(to.path));
  });
}


export default router;
