import HomePage from '@/routes/home/HomePage.vue';
import {useAccountStore} from '@/stores/account';
import {useConfigStore} from '@/stores/config';
import {useUiConfigStore} from '@/stores/ui-config';
import {createRouter, createWebHistory} from 'vue-router';

const routes = [
  {
    path: '/home',
    name: 'home',
    component: HomePage
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/routes/login/LoginPage.vue')
  },
  {
    path: '/setup',
    name: 'setup',
    component: () => import('@/routes/setup/SetupPage.vue')
  },
  // everything else to redirect to the home page
  {
    path: '/:pathMatch(.*)*',
    redirect: '/home'
  }
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
});

if (window) {
  router.beforeEach(async (to, from, next) => {
    const uiConfig = useUiConfigStore();
    await uiConfig.loadConfig();
    const accountStore = useAccountStore();
    // Initialize Local and Keycloak auth instances,
    // so we can check if the user is logged in and/or manage the login flow
    try {
      await accountStore.initialise(uiConfig.config?.auth?.providers);
    } catch (e) {
      console.error('Failed to initialize the account store', e);
    }
    const configStore = useConfigStore();
    const go = (dst) => {
      if (to.path === dst) {
        next();
      } else {
        next(dst);
      }
    };
    const intercept = (dst) => {
      if (to.path !== dst) {
        // Store the current path to redirect back after login
        window.sessionStorage.setItem('redirect', to.fullPath);
      }
      go(dst);
    };

    const authDisabled = uiConfig.auth.disabled;
    const isAuthenticated = accountStore.isLoggedIn;
    const isConfigured = configStore.isConfigured;

    if (!authDisabled && !isAuthenticated) {
      intercept('/login');
      return;
    }

    if (!isConfigured) {
      intercept('/setup');
      return;
    }

    const savedRedirect = window.sessionStorage.getItem('redirect');
    if (savedRedirect) {
      window.sessionStorage.removeItem('redirect');
      go(savedRedirect);
      return;
    }

    next();
  });
}

export default router;
