import auth from '@/routes/auth/route.js';
import devices from '@/routes/devices/route.js';
import ops from '@/routes/ops/route.js';
import automations from '@/routes/automations/route.js';
import site from '@/routes/site/route.js';
import system from '@/routes/system/route.js';
import {route, routeTitle} from '@/util/router.js';
import Vue, {nextTick} from 'vue';
import VueRouter from 'vue-router';
import {useAppConfigStore} from '@/stores/app-config';
import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';

Vue.use(VueRouter);

const router = new VueRouter({
  mode: 'history',
  routes: [
    ...route(auth),
    ...route(devices),
    ...route(ops),
    ...route(automations),
    ...route(site),
    ...route(system)
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
    nextTick(() => (window.document.title = title));
  });
  router.beforeEach(async (to, from, next) => {
    const appConfig = useAppConfigStore();
    const pageStore = storeToRefs(usePageStore());
    await appConfig.loadConfig();
    if (to.path === '/') {
      next(appConfig.homePath);
    } else {
      next(appConfig.pathEnabled(to.path));
    }

    pageStore.showSidebar.value = false;
    pageStore.sidebarTitle.value = '';
    pageStore.sidebarData.value = {};
    pageStore.listedDevice.value = {};
  });
}

export default router;
