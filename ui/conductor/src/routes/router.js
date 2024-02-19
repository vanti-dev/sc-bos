import auth from '@/routes/auth/route.js';
import automations from '@/routes/automations/route.js';
import devices from '@/routes/devices/route.js';
import ops from '@/routes/ops/route.js';
import site from '@/routes/site/route.js';
import system from '@/routes/system/route.js';
import {useAccountStore} from '@/stores/account';
import {usePageStore} from '@/stores/page';
import {useUiConfigStore} from '@/stores/ui-config';
import {route, routeTitle} from '@/util/router.js';

import Vue, {nextTick} from 'vue';
import VueRouter from 'vue-router';

Vue.use(VueRouter);

const router = new VueRouter({
  mode: 'history',
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('./login/LoginPage.vue'),
      meta: {authentication: {rolesRequired: false}, title: 'Login'}
    },
    ...route(auth),
    ...route(devices),
    ...route(ops),
    ...route(automations),
    ...route(site),
    ...route(system)
  ]
});

if (window) {
  router.beforeEach(async (to, from, next) => {
    const uiConfig = useUiConfigStore();
    await uiConfig.loadConfig();
    const authDisabled = uiConfig.config.disableAuthentication;
    const accountStore = useAccountStore();
    // Initialize Local and Keycloak auth instances,
    // so we can check if the user is logged in and/or manage the login flow
    try {
      await accountStore.initialise(uiConfig.config?.auth?.providers);
    } catch (e) {
      console.error('Failed to initialize the account store', e);
    }

    // ------------------------ Data store logic ------------------------ //

    const pageStore = usePageStore();

    // Reset the sidebar to defaults if the path has changed
    if (to.path !== from.path) {
      pageStore.resetSidebarToDefaults();
    }

    // ------------------------ NavigationGuard logic ------------------------ //
    /**
     * Navigation guard for handling route access based on authentication status and application configuration(s).
     *
     * Step 1. Check if the requested path is enabled in the config. If it's not enabled
     *         and the user is trying to access a path other than '/login', redirect to the home path.
     *
     * Step 2. If authentication is not disabled (authDisabled is false):
     *    a. If the user is not logged in:
     *       i.   If the user is trying to access a non-login path, store the current path in session storage
     *            for redirecting back after login, and then redirect the user to the '/login' page.
     *    b. If the user is logged in:
     *       i.   If the user is trying to access the '/login' path, redirect them to the home path.
     *       ii.  If the user is navigating to a regular page (not the '/login' page), check for a saved
     *            redirect path in session storage. If a saved redirect path exists, remove it from session
     *            storage and redirect the user to that path.
     *
     * Finally, if none of the above conditions are met, proceed to the next route (allow the navigation).
     */
    const isPathEnabled = uiConfig.pathEnabled(to.path);
    const redirectToHome = () => next(uiConfig.homePath);
    const isLoginPath = to.path === '/login';
    const isAuthenticated = accountStore.isLoggedIn;

    if (!isPathEnabled && (!isLoginPath || authDisabled && isLoginPath && from.path !== uiConfig.homePath)) {
      redirectToHome();
      return;
    }

    if (!authDisabled) {
      if (!isAuthenticated) {
        if (to.path !== '/login') {
          // Store the current path to redirect back after login
          window.sessionStorage.setItem('redirect', to.fullPath);
          next('/login');
          return;
        }
      } else {
        if (to.path === '/login') {
          // Redirect logged-in users away from the login page to home
          redirectToHome();
          return;
        } else {
          // If navigating to a regular page, check for a saved redirect path
          const savedRedirect = window.sessionStorage.getItem('redirect');
          if (savedRedirect) {
            window.sessionStorage.removeItem('redirect');
            next(savedRedirect); // Redirect to the saved path
            return;
          }
        }
      }
    }

    next();
  });

  router.afterEach((to, from) => {
    const nt = routeTitle(to);
    const ot = routeTitle(from);
    if (nt === ot) {
      return;
    }

    const title = nt ? `${nt} - Smart Core` : `Smart Core`;
    nextTick(() => window.document.title = title);
  });
}

export default router;
