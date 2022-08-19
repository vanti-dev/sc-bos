import Keycloak from 'keycloak-js';

/**
 * @type {Promise<Keycloak> | null}
 */
let instance = null;
/**
 * Events from the keycloak instance.
 *
 * Supported events are and follow the semantics and payload of the corresponding onXxx keycloak methods:
 *  - ready
 *  - authError
 *  - authLogout
 *  - authSuccess
 *  - authRefreshError
 *  - authRefreshSuccess
 *  - tokenExpired
 *
 * @type {EventTarget}
 */
export const events = new EventTarget();

/**
 * @returns {Promise<Keycloak>}
 */
export function keycloak() {
  if (instance === null) {
    instance = newKeycloak();
  }
  return instance;
}

/**
 * @returns {Promise<Keycloak>}
 */
async function newKeycloak() {
  const kc = new Keycloak(await constructorConfig());
  // setup event handling
  kc.onReady = authenticated => {
    const e = new Event('ready');
    e.authenticated = authenticated;
    events.dispatchEvent(e)
  }
  kc.onAuthError = errorData => {
    const e = new Event('authError');
    e.errorData = errorData;
    events.dispatchEvent(e);
  }
  kc.onAuthLogout = () => events.dispatchEvent(new Event('authLogout'));
  kc.onAuthSuccess = () => events.dispatchEvent(new Event('authSuccess'));
  kc.onAuthRefreshError = () => events.dispatchEvent(new Event('authRefreshError'));
  kc.onAuthRefreshSuccess = () => events.dispatchEvent(new Event('authRefreshSuccess'));
  kc.onTokenExpired = () => events.dispatchEvent(new Event('tokenExpired'));

  await kc.init(await initConfig());
  return kc;
}

/**
 * @returns {Promise<import('keycloak-js').KeycloakConfig | string>}
 */
async function constructorConfig() {
  // todo: get keycloak config from somewhere non-hard-coded
  return {
    realm: 'smart-core',
    url: 'http://localhost:8888/',
    clientId: 'sc-apps'
  }
}

/**
 * @returns {Promise<import('keycloak-js').KeycloakInitOptions>}
 */
async function initConfig() {
  // todo: get keycloak init config from somewhere non-hard-coded
  return {}
}
