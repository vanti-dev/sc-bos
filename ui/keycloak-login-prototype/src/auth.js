import Keycloak from 'keycloak-js'
import {computed, ref} from 'vue';

export const keycloakClientId = 'scos-opsui'
export const keycloakBaseURL = 'http://localhost:8888/';
export const keycloakRealm = 'Smart_Core';

export const keycloak = new Keycloak({
  realm: keycloakRealm,
  url: keycloakBaseURL,
  clientId: keycloakClientId
})
export const openIDConfigURL = `${keycloakBaseURL}/realms/${keycloakRealm}/.well-known/openid-configuration`

let oidcConfigP = null

export function useOpenIDConnect(configURL = openIDConfigURL) {
  const config = ref(null)
  const configError = ref(null)

  if (oidcConfigP === null) {
    oidcConfigP = fetch(configURL)
        .then(res => {
          if (!res.ok) {
            throw new Error(`Failed to get OpenID Connect configuration from ${configURL}`)
          }
          return res.json()
        })
  }

  // this is done per call to useOpenIDConnect, but we only fetch the data once
  oidcConfigP
      .then(cfg => {
        console.debug('OpenID Connect configuration', cfg);
        configError.value = null
        config.value = cfg
      })
      .catch(err => {
        configError.value = err
        config.value = null
      })

  const tokenEndpoint = computed(() => {
    return config.value?.token_endpoint ?? null
  });
  const deviceAuthorizationEndpoint = computed(() => {
    return config.value?.device_authorization_endpoint ?? null
  });
  return {
    config, configError,
    allScopes: computed(() => {
      return config.value?.scopes_supported?.sort() ?? []
    }),

    // general OAuth props and functions
    tokenEndpoint,
    doRefreshToken(refreshToken) {
      return oidcConfigP
          .then(() => tokenEndpoint.value)
          .then(url => {
            return fetch(url, {
              method: 'POST',
              body: new URLSearchParams({
                client_id: keycloakClientId,
                grant_type: 'refresh_token',
                refresh_token: refreshToken
              })
            })
          })
          .then(res => {
            if (!res.ok) {
              throw new Error(`Failed to refresh token: ${res.status} ${res.statusText}`)
            }
            return res.json()
          });
    },

    // OAuth Device Auth Flow
    deviceAuthorizationEndpoint,
    beginDeviceAuth(scopes = ['profile']) {
      return oidcConfigP
          .then(() => deviceAuthorizationEndpoint.value)
          .then(url => {
            return fetch(url, {
              method: 'POST',
              body: new URLSearchParams({
                client_id: keycloakClientId,
                scope: scopes.join(' ')
              })
            })
          })
          .then(res => {
            if (!res.ok) {
              throw new Error(`Failed to start device auth flow: ${res.status} ${res.statusText}`)
            }
            return res.json()
          });
    },
    checkDeviceToken(lastResponse) {
      return oidcConfigP
          .then(() => tokenEndpoint.value)
          .then(url => {
            return fetch(url, {
              method: 'POST',
              body: new URLSearchParams({
                client_id: keycloakClientId,
                grant_type: 'urn:ietf:params:oauth:grant-type:device_code',
                device_code: lastResponse.device_code
              })
            });
          })
          .then(res => {
            if (!res.ok) {
              throw new Error(`Failed to check device auth flow: ${res.status} ${res.statusText}`)
            }
            return res.json()
          });
    }
  }
}

export function decodeToken(token) {
  return JSON.parse(atob(token.split('.')[1]))
}

