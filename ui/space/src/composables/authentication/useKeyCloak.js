import {events, keycloak} from '@/api/keycloak.js';

/**
 * Keycloak composable
 * This composable is responsible for handling all Keycloak related functionality
 *
 * @return {AuthenticationProvider & {
 *  kcp: Promise<import('keycloak-js').KeycloakInstance>,
 *  kcEvents: import('keycloak-js').KeycloakEventEmitter,
 *  login: (function(import('keycloak-js').KeycloakLoginOptions?): Promise<AuthenticationDetails>),
 * }}
 */
export default function() {
  const kcp = keycloak();
  const kcEvents = events;


  /**
   * Convert a keycloak instance to AuthenticationDetails.
   *
   * @param {import('keycloak-js').Keycloak} kc
   * @return {Promise<AuthenticationDetails>}
   */
  const kcToAuthDetails = (kc) => {
    return {
      claims: kc.idTokenParsed,
      loggedIn: kc.authenticated,
      token: kc.token
    };
  };
  //
  // ----------------------------------- //
  //
  /**
   * Initialize the Keycloak instance and event listeners
   *
   * @return {Promise<AuthenticationDetails|null>}
   */
  const initializeKeycloak = async () => {
    const kc = await kcp;
    if (!kc.authenticated) {
      return null;
    }
    return kcToAuthDetails(kc);
  };

  /**
   * Login to Keycloak with the given scopes
   *
   * @param {import('keycloak-js').KeycloakLoginOptions} [options]
   * @return {Promise<AuthenticationDetails>}
   */
  const loginKeyCloak = async (options) => {
    const kc = await kcp;
    await kc.login(options);
    // not needed as login will redirect the page, but helpful for js type checking
    return undefined;
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
   * @param {number} [leeway]
   * @return {Promise<AuthenticationDetails>}
   */
  const refreshToken = async (leeway = 15) => {
    const kc = await kcp;
    await kc.updateToken(leeway);
    return kcToAuthDetails(kc);
  };

  return {
    kcp,
    kcEvents,

    init: initializeKeycloak,
    login: loginKeyCloak,
    logout: logoutKeyCloak,
    refreshToken
  };
}
