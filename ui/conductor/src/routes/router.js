import auth from '@/routes/auth/route.js';
import devices from '@/routes/devices/route.js';
import ops from '@/routes/ops/route.js';
import automations from '@/routes/automations/route.js';
import site from '@/routes/site/route.js';
import system from '@/routes/system/route.js';
import {useAccountStore} from '@/stores/account';
import {route, routeTitle} from '@/util/router.js';
import Vue, {computed, nextTick} from 'vue';
import VueRouter from 'vue-router';
import {useAppConfigStore} from '@/stores/app-config';
import useAuthSetup from '@/composables/useAuthSetup';
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
    nextTick(() => window.document.title = title);
  });
  router.beforeEach(async (to, from, next) => {
    const appConfig = useAppConfigStore();
    const {loginDialog, isLoggedIn} = storeToRefs(useAccountStore());
    await appConfig.loadConfig();

    // Update the meta tag based on disableAuthentication
    // This way we can use the meta tag in the router to determine if authentication is required on a page
    if (appConfig.config) {
      router.getRoutes().forEach(route => {
        // Update the meta tag based on disableAuthentication
        if (route.meta.requiresAuth !== undefined) {
          route.meta.requiresAuth = !appConfig.config.disableAuthentication;
        }
      });
    }

    // Clear the sidebar when navigating to a different main path
    const {showSidebar, sidebarTitle, sidebarData} = storeToRefs(usePageStore());
    const mainPathFrom = from.path.split('/')[1];
    const mainPathTo = to.path.split('/')[1];

    if (mainPathFrom !== mainPathTo) {
      showSidebar.value = false;
      sidebarTitle.value = '';
      sidebarData.value = {};
    }

    const {hasNoAccess} = useAuthSetup();
    const userLoggedIn = computed(() => appConfig.config.disableAuthentication || isLoggedIn.value);
    const isPathEnabled = computed(() => appConfig.pathEnabled(to.path));
    const authenticationRequired = to.meta?.requiresAuth;
    // Display login modal instantly on screen if user is not logged in,
    // but we require authentication on the page we visit
    loginDialog.value = authenticationRequired;


    /**
     * Handles route navigation based on several conditions:
     *
     * 1. If the path is enabled, we can try navigating to it
     *   - If the path requires authentication and the user is not logged in, we display the login dialog
     *   - If the path is not accessible, we navigate to the home path
     *   - Otherwise, we navigate to the path
     * 2. If the path is disabled, we navigate to the home path
     */
    if (isPathEnabled.value) {
      if (authenticationRequired && !userLoggedIn.value) {
        loginDialog.value = true;
        next();
      } else if (hasNoAccess(to.path)) {
        next(appConfig.homePath);
      } else {
        next();
      }
    } else {
      next(appConfig.homePath);
    }
  });
}

export default router;
