import Keycloak from 'keycloak-js'

export const keycloakClientId = 'scos-opsui'

export const keycloak = new Keycloak({
  realm: 'Smart_Core',
  url: 'http://localhost:8888/',
  clientId: keycloakClientId
})

export const allScopes = [
  'profile',
  'Read',
  'Write'
]
