import {localLogin} from '@/api/localLogin.js';
import jwtDecode from 'jwt-decode';
import {useAccountStore} from '@/stores/account';

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

        if (payload.access_token) {
          accountStore.authenticationDetails.claims = {
            email: username,
            ...jwtDecode(payload.access_token)
          };
          accountStore.authenticationDetails.loggedIn = !!payload.access_token;
          accountStore.authenticationDetails.token = payload.access_token;
        }
      } else {
        accountStore.snackbar = true;
      }
    } catch (err) {
      accountStore.snackbar = true;
    }
  };

  /**
   * Logout using the store reset function
   *
   * @return {void}
   */
  const logoutLocal = () => {
    accountStore.resetStoreToDefaults();
  };

  return {
    login: loginLocal,
    logout: logoutLocal
  };
}
