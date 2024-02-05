import {localLogin} from '@/api/localLogin.js';
import {useAccountStore} from '@/stores/account';
import {loadFromBrowserStorage, saveToBrowserStorage} from '@/util/browserStorage';
import jwtDecode from 'jwt-decode';
import {ref} from 'vue';

/**
 * Local authentication composable
 * This composable is responsible for handling all local authentication related functionality
 *
 * @return {{
 *  login: (function(string, string): Promise<void>),
 *  logout: (function(): void)
 * }}
 */
export default function() {
  const accountStore = useAccountStore();

  // ----------------- //

  /**
   * Load the local authentication details from local storage
   *
   * @type {import('vue').Ref<import('@/stores/account').AuthenticationDetails|null>}
   */
  const existingLocalAuth = ref(null);

  /**
   * Initialize local authentication - set the existing local authentication details
   * or the default store values
   *
   * @return {import('@/stores/account').AuthenticationDetails|null}
   */
  const initializeLocal = () => {
    existingLocalAuth.value = loadFromBrowserStorage(
        'local',
        'authenticationDetails',
        accountStore.authenticationDetails
    )[0];

    return existingLocalAuth.value;
  };

  /**
   * Login using local authentication provider
   *
   * @param {string} username
   * @param {string} password
   * @return {Promise<void>}
   */
  const loginLocal = async (username, password) => {
    try {
      const res = await localLogin(username, password);

      if (res.status === 200) {
        const payload = await res.json();

        if (payload?.access_token) {
          accountStore.authenticationDetails.claims = {
            email: username,
            ...jwtDecode(payload.access_token)
          };
          accountStore.authenticationDetails.loggedIn = !!payload.access_token;
          accountStore.authenticationDetails.token = payload.access_token;

          accountStore.snackbar = {
            message: 'Failed to sign in, please try again.',
            visible: false
          };
        }
      } else {
        accountStore.snackbar = {
          visible: true,
          message: 'Failed to sign in, please try again.'
        };
      }

      saveToBrowserStorage('local', 'authenticationDetails', accountStore.authenticationDetails);
    } catch (err) {
      accountStore.snackbar = {
        visible: true,
        message: 'Failed to sign in, please try again.'
      };
      saveToBrowserStorage('local', 'authenticationDetails', accountStore.authenticationDetails);
    }
  };

  /**
   * Logout using the store reset function, then store the cleared store in local storage
   *
   * @return {void}
   */
  const logoutLocal = async () => {
    await window.localStorage.removeItem('authenticationDetails');
  };

  return {
    existingLocalAuth,
    init: initializeLocal,
    login: loginLocal,
    logout: logoutLocal
  };
}
