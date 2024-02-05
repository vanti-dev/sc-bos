import useKeyCloak from '@/composables/authentication/useKeyCloak';
import useLocal from '@/composables/authentication/useLocal';
import {useAppConfigStore} from '@/stores/app-config';
import {loadFromBrowserStorage} from '@/util/browserStorage';
import {defineStore} from 'pinia';
import {computed, ref, watch} from 'vue';
import {useRouter} from 'vue-router/composables';

/**
 * @typedef AuthenticationDetails
 * @property {Object<string,*>} claims
 * @property {boolean} loggedIn
 * @property {string} token
 */

export const useAccountStore = defineStore('accountStore', () => {
  const appConfig = useAppConfigStore();
  const keyCloak = useKeyCloak();
  const localAuth = useLocal();
  const router = useRouter();

  // Set up the storage for the login: authProvider, claims, login status and token
  const authenticationDetails = ref(
      /** @type {AuthenticationDetails & {authProvider: string}} */
      {
        authProvider: '',
        claims: {},
        loggedIn: false,
        token: ''
      });
  const loginFormVisible = ref(false);
  const snackbar = ref({
    message: 'Failed to sign in, please try again',
    visible: false
  });


  /**
   * Reset the store values to defaults
   *
   * @return {void}
   */
  const resetStoreToDefaults = () => {
    authenticationDetails.value = {
      authProvider: '',
      claims: {},
      loggedIn: false,
      token: ''
    };
  };


  /**
   * Initialize Keycloak and Local Auth instances, so we can check if the user is logged in and/or manage the login flow
   *
   * @return {Promise<void>}
   */
  const initialise = async () => {
    if (appConfig.config.disableAuthentication) {
      return;
    }
    try {
      // Attempt to initialize Keycloak authentication
      if (appConfig.config?.keycloak) {
        const kcResponse = await keyCloak.init();
        if (kcResponse) {
          availableAuthProviders.value = ['keyCloakAuth', 'localAuth'];
        } else {
          availableAuthProviders.value = ['localAuth'];
          return;
        }
        if (kcResponse?.authenticated) {
          authenticationDetails.value.authProvider = 'keyCloakAuth';
          return; // Exit if authenticated with Keycloak
        }
      }
    } catch (error) {
      console.error('Keycloak initialization failed', error);
      snackbar.value = {
        message: 'Keycloak initialization failed: ' + error.error,
        visible: true
      };
      // Proceed to displaying the local authentication form if Keycloak fails
      loginFormVisible.value = true;
    }
    // Initialize local authentication if Keycloak is not configured, fails, or is not authenticated
    try {
      authenticationDetails.value = await localAuth.init();
      if (authenticationDetails.value.loggedIn) {
        return; // Exit if authenticated with local auth
      }
    } catch (error) {
      console.error('Local authentication initialization failed', error);
      snackbar.value = {
        message: 'Local authentication initialization failed: ' + error.error,
        visible: true
      };
      resetStoreToDefaults();
    }
  };
  //
  // ----------------------------------- //
  //
  /**
   * Check if authentication is disabled
   *
   * @type {import('vue').ComputedRef<boolean>}
   */
  const isAuthenticationDisabled = computed(() => {
    return appConfig.config?.disableAuthentication || false;
  });

  /**
   * Collect the available authentication providers
   *
   * @type {import('vue').Ref<string[]>}
   */
  const availableAuthProviders = ref(['localAuth']);

  /**
   * Retrieve the authentication provider used for login
   *
   * @type {import('vue').ComputedRef<string|null>}
   */
  const activeAuthProvider = computed(() => {
    return authenticationDetails.value.authProvider || null;
  });

  /**
   * Set the authentication provider depending on the login form visibility
   * If the login form is visible, use the local authentication provider, otherwise use KeyCloak
   */
  const watchSources = [availableAuthProviders, isAuthenticationDisabled];

  watch(watchSources, ([availableProviders, authDisabled]) => {
    if (!authDisabled) {
      loginFormVisible.value = !availableProviders.includes('keyCloakAuth');
    }
  }, {immediate: true, deep: true});

  /**
   * Returns the login status depending on the authentication provider
   *
   * @type {import('vue').ComputedRef<boolean>}
   */
  const isLoggedIn = computed(() => {
    const detailsAvailable = !!(
      authenticationDetails.value.claims &&
        authenticationDetails.value.loggedIn &&
        authenticationDetails.value.token
    );

    return detailsAvailable || isAuthenticationDisabled.value;
  });

  /**
   * Returns the full name of the logged in user
   *
   * @type {import('vue').ComputedRef<string>}
   */
  const fullName = computed(() => authenticationDetails.value.claims?.name || '');

  /**
   * Returns the email of the logged in user
   *
   * @type {import('vue').ComputedRef<string>}
   */
  const email = computed(() => authenticationDetails.value.claims?.email || '');

  /**
   * Returns the roles of the logged in user
   *
   * @type {import('vue').ComputedRef<string[]>}
   */
  const roles = computed(() => authenticationDetails.value.claims?.roles || []);
  //
  // ----------------------------------- //
  //
  // Dynamic controls - depending on the active authentication provider
  //
  /** @typedef {{ username: string, password: string }} LocalAuthLoginValues */

  /**
   * Log in with local authentication using the given username and password
   *
   * @param {LocalAuthLoginValues} values
   * @return {Promise<void>}
   */
  const loginWithLocalAuth = async (values) => {
    authenticationDetails.value.authProvider = 'localAuth';
    await localAuth.login(values.username, values.password);

    // If the login was successful
    if (isLoggedIn.value) {
      // If there is a redirect in the session storage, redirect to that page
      const redirect = loadFromBrowserStorage('session', 'redirect', '')[0];
      if (redirect !== '') {
        window.sessionStorage.removeItem('redirect');
        await router.push(redirect);

        // Otherwise, redirect to the home page
      } else {
        await router.push(appConfig.homePath);
      }
    }
  };

  /**
   * Log in with KeyCloak using the given scopes
   *
   * @param {string[]} scopes
   * @return {Promise<void>}
   */
  const loginWithKeyCloak = async (scopes) => {
    await keyCloak.login(scopes);
  };

  /**
   * Logout of the active authentication provider, then reset the store to defaults.
   * If a reason is provided, display a snackbar with the reason.
   *
   * @param {string} reason
   * @return {Promise<void>}
   */
  const logout = async (reason) => {
    if (activeAuthProvider.value === 'keyCloakAuth') {
      await keyCloak.logout();
    } else if (activeAuthProvider.value === 'localAuth') {
      await localAuth.logout();
    }

    if (reason) {
      snackbar.value = {
        message: 'Logged out: ' + reason,
        visible: true
      };
    }

    resetStoreToDefaults();
    window.localStorage.removeItem('authenticationDetails');

    if (router.currentRoute.path !== '/login') {
      await router.push('/login');
    }
  };

  const refreshToken = async () => {
    if (activeAuthProvider.value === 'keyCloakAuth') {
      await keyCloak.refreshToken();
    }
  };

  return {
    initialise,
    authenticationDetails,
    loginFormVisible,
    snackbar,
    resetStoreToDefaults,

    isAuthenticationDisabled,
    availableAuthProviders,
    isLoggedIn,
    fullName,
    email,
    roles,

    loginWithLocalAuth,
    loginWithKeyCloak,
    logout,
    refreshToken
  };
});
