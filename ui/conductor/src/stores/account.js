import useKeyCloak from '@/composables/authentication/useKeyCloak';
import useLocal from '@/composables/authentication/useLocal';
import {useAppConfigStore} from '@/stores/app-config';
import {loadFromBrowserStorage} from '@/util/browserStorage.js';
import {defineStore} from 'pinia';
import {computed, ref, watch} from 'vue';

export const useAccountStore = defineStore('accountStore', () => {
  const appConfig = useAppConfigStore();
  const keyCloak = useKeyCloak();
  const localAuth = useLocal();

  // Set up the storage for the login: authProvider, claims, login status and token
  const authenticationDetails = ref({
    authProvider: '',
    claims: {},
    loggedIn: false,
    token: ''
  });
  const loginFormVisible = ref(false);
  const redirect = ref(...loadFromBrowserStorage('session', 'redirect', ''));
  const snackbar = ref(false);

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
  //
  // ----------------------------------- //
  //
  /**
   * Check if authentication is disabled
   *
   * @type {import('vue').ComputedRef<boolean>}
   */
  const isAuthenticationDisabled = computed(() => {
    return appConfig.config.disableAuthentication;
  });

  /**
   * Collect the available authentication providers
   *
   * @type {import('vue').ComputedRef<string[]>}
   */
  const availableAuthProviders = computed(() => {
    if (appConfig.config.keycloak) {
      return ['keyCloakAuth', 'localAuth'];
    } else {
      return ['localAuth'];
    }
  });

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
      authenticationDetails.value.authProvider &&
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

  /**
   * Returns the redirect path if it exists
   *
   * @type {import('vue').ComputedRef<string|null>}
   */
  const activeRedirect = computed(() => redirect.value || null);
  //
  // ----------------------------------- //
  //
  // Dynamic controls - depending on the active authentication provider
  //
  /** @typedef {{ username: string, password: string }} LocalAuthLoginValues */

  /**
   * Login to the active authentication provider
   *
   * @param {LocalAuthLoginValues|string[]} values
   */
  const loginWithLocalAuth = async (values) => {
    await localAuth.login(values.username, values.password);
  };

  const loginWithKeyCloak = async (scopes) => {
    await keyCloak.login(scopes);
  };

  /**
   * Logout of the active authentication provider, then reset the store to defaults
   *
   * @return {Promise<void>}
   */
  const logout = async () => {
    if (activeAuthProvider.value === 'keyCloakAuth') {
      await keyCloak.logout();
    } else {
      localAuth.logout();
    }

    resetStoreToDefaults();
  };

  const refreshToken = async () => {
    if (activeAuthProvider.value === 'keyCloakAuth') {
      await keyCloak.refreshToken();
    }
    // else {
    //   localAuth.refreshToken();
    // }
  };

  return {
    authenticationDetails,
    loginFormVisible,
    snackbar,
    resetStoreToDefaults,

    isAuthenticationDisabled,
    availableAuthProviders,
    activeAuthProvider,
    activeRedirect,
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
