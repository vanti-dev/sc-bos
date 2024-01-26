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
import useAuthSetup from '@/composables/useAuthSetup';
import {usePageStore} from '@/stores/page';
import {useOverviewStore} from '@/routes/ops/overview/overviewStore.js';
import {storeToRefs} from 'pinia';

Vue.use(VueRouter);

const router = new VueRouter({
  mode: 'history',
  routes: [
    // {
    //   path: '/login',
    //   name: 'login',
    //   component: () => import('./login/LoginPage.vue'),
    //   meta: {authentication: {rolesRequired: false}, title: 'Login'}
    // },
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
    await appConfig.loadConfig();

    const authSetup = useAuthSetup();
    authSetup.navigate(to.path, next);

    // Clear the sidebar when navigating to a different main path
    const {showSidebar, sidebarTitle, sidebarData} = storeToRefs(usePageStore());
    const mainPathFrom = from.path.split('/')[1];
    const mainPathTo = to.path.split('/')[1];

    if (mainPathFrom !== mainPathTo) {
      showSidebar.value = false;
      sidebarTitle.value = '';
      sidebarData.value = {};
    }

    const {activeOverview} = storeToRefs(useOverviewStore());

    /**
     * Select the relevant child item when navigating to a building overview page with a child item
     * For example, when navigating to '/ops/overview/building/Floor%203/Left%20Wing'
     * This helps when the user refreshes the page or navigates directly to the page (shared link)
     */
    const buildingChildren = appConfig.config?.building?.children;
    const currentPathSegments = to.path.split('/').filter(segment => segment);
    const lastSegment = currentPathSegments[currentPathSegments.length - 1];
    const findItemByTitle = (items, title) => {
      for (const item of items) {
        if (encodeURIComponent(item.title) === title) {
          return item;
        }
        if (item.children) {
          const found = findItemByTitle(item.children, title);
          if (found) return found;
        }
      }
      return null;
    };

    // Find the active item in the building children
    const activeItem = findItemByTitle(buildingChildren, lastSegment);

    if (activeItem) {
      // eslint-disable-next-line no-unused-vars
      const {children, ...item} = activeItem;

      activeOverview.value = item;
    }

    // Clear the activeOverview when navigating to a different page - other than overview
    // Check if the path is not '/ops/overview/building' or doesn't start with '/ops/overview/building'
    if (!to.path.startsWith('/ops/overview/building')) {
      // Clear the activeOverview
      activeOverview.value = null;
    }
  });
}

export default router;
