import Keycloak from "keycloak-js"

export const keycloakClientId = "sc-apps"

export const keycloak = new Keycloak({
  realm: "smart-core",
  url: "http://localhost:8888/",
  clientId: keycloakClientId
})

export const allScopes = [
  "Test.Read",
  "Test.Write"
]