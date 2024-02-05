import {events, keycloak} from '@/api/keycloak.js';
import {useAccountStore} from '@/stores/account';

/**
 * Keycloak composable
 * This composable is responsible for handling all Keycloak related functionality
 *
 * @return {{
 *  kcp: Promise<import('keycloak-js').KeycloakInstance>,
 *  kcEvents: import('keycloak-js').KeycloakEventEmitter,
 *  initializeKeycloak: (function(): Promise<void>),
 *  login: (function(*=): Promise<void>),
 *  logout: (function(): Promise<void>)
 * }}
 */
export default function() {
  const accountStore = useAccountStore();
  const kcp = keycloak();
  const kcEvents = events;


  /**
   * Update the auth status and save to local storage
   *
   * @return {Promise<void>}
   */
  const updateAuthStatus = async () => {
    const kc = await kcp;
    accountStore.authenticationDetails.claims = kc.idTokenParsed;
    accountStore.authenticationDetails.loggedIn = kc.authenticated;
    accountStore.authenticationDetails.token = kc.token;
  };
  //
  // ----------------------------------- //
  //
  /**
   * Initialize the Keycloak instance and event listeners
   *
   * @return {Promise<import('keycloak-js').KeycloakInstance>}
   */
  const initializeKeycloak = async () => {
    const kc = await kcp;
    if (kc.authenticated) {
      await updateAuthStatus();
    }

    return kc;
  };

  /**
   * Login to Keycloak with the given scopes
   *
   * @param {string[]} scopes
   * @return {Promise<void>}
   */
  const loginKeyCloak = async (scopes) => {
    const kc = await kcp;
    kc.login({scope: scopes.join(' ')});
  };

  /**
   * Logout of Keycloak and clear all login details
   *
   * @return {Promise<void>}
   */
  const logoutKeyCloak = async () => {
    const kc = await kcp;
    kc.logout();
  };

  /**
   * Update the token if it is close to expiring (15 seconds)
   *
   * @return {Promise<void>}
   */
  const refreshToken = async () => {
    const kc = await kcp;
    await kc.updateToken(15);
  };
  //
  // ----------------------------------- //
  //
  /**
   * AuthSuccess event listener
   * This event listener is responsible for updating the authenticationDetails on a successful login
   */
  kcEvents.addEventListener('authSuccess', async () => await updateAuthStatus());

  /**
   * AuthRefreshSuccess event listener
   * This event listener is responsible for updating the authenticationDetails on a successful token refresh
   */
  kcEvents.addEventListener('onAuthRefreshSuccess', async () => await updateAuthStatus());

  return {
    kcp,
    kcEvents,

    init: initializeKeycloak,
    login: loginKeyCloak,
    logout: logoutKeyCloak,
    refreshToken
  };
}
