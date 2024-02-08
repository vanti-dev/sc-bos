import useKeyCloak from '@/composables/authentication/useKeyCloak';
import useLocal from '@/composables/authentication/useLocal';
import {useUiConfigStore} from '@/stores/ui-config';
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
/**
 * Describes the interface an authentication provider should implement.
 * Providers should also provide a mechanism for beginning the authentication flow.
 *
 * @typedef AuthenticationProvider
 * @property {function(): Promise<AuthenticationDetails|null>} init
 *  Initialise the auth provider, returning null if not authenticated
 * @property {function(): Promise<void>} logout
 * @property {function(): Promise<AuthenticationDetails|null>} refreshToken
 */

export const useAccountStore = defineStore('accountStore', () => {
  const uiConfig = useUiConfigStore();
  const keyCloak = useKeyCloak();
  const localAuth = useLocal();
  const router = useRouter();

  // initComplete is resolved (or rejected) the first time initialise is called.
  // Functions can `await initComplete` to make sure that any authenticationDetails -
  // including authProvider - are set correctly.
  let initResolved;
  let initRejected;
  const initComplete = new Promise((resolve, reject) => {
    initResolved = resolve;
    initRejected = reject;
  });

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
   * Helper for initialise that contains all the logic so any 'finally' logic can be run for all cases.
   *
   * @return {Promise<void>}
   */
  const _initialise = async () => {
    if (uiConfig.config.disableAuthentication) {
      return;
    }

    const providers = [
      {
        name: 'keyCloakAuth',
        init: keyCloak.init,
        enabled: () => Boolean(uiConfig.config?.keycloak)
      },
      {
        name: 'localAuth',
        init: localAuth.init,
        enabled: () => true
      }
    ];

    let loginDetails = null; // details from the first successful init attempt (that returned 'logged in')
    const availableProviderNames = [];
    for (const provider of providers) {
      if (!provider.enabled()) continue;
      try {
        if (loginDetails === null) {
          const details = await provider.init();
          if (details) {
            loginDetails = {
              ...details,
              authProvider: provider.name
            };
          }
        }
        availableProviderNames.push(provider.name);
      } catch (e) {
        console.error(`${provider.name} initialization failed`, e);
        snackbar.value = {
          message: `${provider.name} initialization failed: ${e.error}`,
          visible: true
        };
      }
    }

    availableAuthProviders.value = availableProviderNames;
    if (loginDetails) {
      // we are logged in already
      authenticationDetails.value = loginDetails;
      return;
    }

    // not logged in, show the login form
    loginFormVisible.value = true;
  };

  /**
   * Initialize Keycloak and Local Auth instances, so we can check if the user is logged in and/or manage the login flow
   *
   * @return {Promise<void>}
   */
  const initialise = async () => {
    try {
      await _initialise();
      initResolved();
    } catch (e) {
      initRejected(e);
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
    return uiConfig.config?.disableAuthentication || false;
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
    const detailsAvailable = !!(authenticationDetails.value.token);

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
   * Perform a login using the given fn and store details and redirect if needed.
   *
   * @param {string} authProvider
   * @param {function(): Promise<AuthenticationDetails>} fn
   * @return {Promise<void>}
   */
  const doLogin = async (authProvider, fn) => {
    await initComplete;
    const details = await fn();
    if (details) {
      authenticationDetails.value = {...details, authProvider};
      await redirectToLastPage();
    }
  };

  /**
   * Log in with local authentication using the given username and password
   *
   * @param {LocalAuthLoginValues} values
   * @return {Promise<void>}
   */
  const loginWithLocalAuth = async (values) => {
    return doLogin('localAuth', () => localAuth.login(values.username, values.password));
  };

  /**
   * Log in with KeyCloak using the given scopes
   *
   * @param {string[]} scopes
   * @return {Promise<void>}
   */
  const loginWithKeyCloak = async (scopes) => {
    return doLogin('keyCloakAuth', () => keyCloak.login(scopes));
  };

  /**
   * Redirect to the login page if the user is not already there.
   *
   * @return {Promise<void>}
   */
  const redirectToLogin = async () => {
    if (router.currentRoute.path !== '/login') {
      await router.push('/login');
    }
  };

  /**
   * Redirect to the last page the user was on, or the home page if not set.
   *
   * @return {Promise<void>}
   */
  const redirectToLastPage = async () => {
    // If there is a redirect in the session storage, redirect to that page
    const redirect = loadFromBrowserStorage('session', 'redirect', '')[0];
    if (redirect !== '') {
      window.sessionStorage.removeItem('redirect');
      await router.push(redirect);

      // Otherwise, redirect to the home page
    } else {
      await router.push(uiConfig.homePath);
    }
  };

  /**
   * Logout of the active authentication provider, then reset the store to defaults.
   * If a reason is provided, display a snackbar with the reason.
   *
   * @param {string} reason
   * @return {Promise<void>}
   */
  const logout = async (reason) => {
    await initComplete;
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
    await redirectToLogin();
  };

  /**
   * @param {function(): Promise<AuthenticationDetails>} fn
   * @return {Promise<void>}
   */
  const doRefreshToken = async (fn) => {
    try {
      const details = await fn();
      authenticationDetails.value = {...details, authProvider: authenticationDetails.value.authProvider};
    } catch (e) {
      resetStoreToDefaults();
      snackbar.value = {
        message: 'Session expired, please log in again: ' + e,
        visible: true
      };
      await redirectToLogin();
    }
  };

  const refreshToken = async () => {
    await initComplete;
    if (activeAuthProvider.value === 'keyCloakAuth') {
      return doRefreshToken(keyCloak.refreshToken);
    } else if (activeAuthProvider.value === 'localAuth') {
      return doRefreshToken(localAuth.refreshToken);
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
