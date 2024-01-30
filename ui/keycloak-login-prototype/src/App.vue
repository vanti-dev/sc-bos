<script setup>
import PreJson from '@/component/PreJson.vue';
import {ref} from 'vue'
import {decodeToken, keycloak, useOpenIDConnect} from './auth'
import GetTest from './component/GetTest.vue'
import UpdateTest from './component/UpdateTest.vue';

const selectedScopes = ref(['profile'])
const accessToken = ref(/** @type {null|string} */ null)
const refreshToken = ref(/** @type {null|string} */ null)
const keycloakClaims = ref(/** @type {null|KeycloakTokenParsed} */ null)

keycloak.init({})
    .then(authenticated => {
      if (authenticated) {
        keycloakClaims.value = keycloak.idTokenParsed
        accessToken.value = keycloak.token
        refreshToken.value = keycloak.refreshToken
      }
    })
    .catch(console.error)

async function loginKeycloak() {
  const scopes = selectedScopes.value.join(' ')

  try {
    await keycloak.login({
      scope: scopes
    })
  } catch (e) {
    console.error(e)
  }
}

const oidc = useOpenIDConnect();
const beginDeviceResponse = ref(null);

async function beginDeviceAuth() {
  beginDeviceResponse.value = await oidc.beginDeviceAuth(selectedScopes.value);
}

const checkDeviceResponse = ref(null);
const checkDeviceError = ref(null);

async function checkDeviceAuth() {
  const lastResponse = checkDeviceResponse.value ?? beginDeviceResponse.value;
  if (!lastResponse) return;
  try {
    checkDeviceResponse.value = await oidc.checkDeviceToken(lastResponse);
    checkDeviceError.value = null;

    accessToken.value = checkDeviceResponse.value.access_token;
    refreshToken.value = checkDeviceResponse.value.refresh_token;
    keycloakClaims.value = decodeToken(accessToken.value);

  } catch (e) {
    checkDeviceError.value = e;
    checkDeviceResponse.value = null;
  }
}

const refreshTokenResponse = ref(null);

async function doRefreshToken() {
  const response = await oidc.doRefreshToken(refreshToken.value);
  refreshTokenResponse.value = response;
  accessToken.value = response.access_token;
  refreshToken.value = response.refresh_token;
  keycloakClaims.value = decodeToken(accessToken.value);
}

async function logoutKeycloak() {
  return keycloak.logout();
}
</script>

<template>
  <h1>Auth Demo</h1>
  <div>
    <h2>Keycloak</h2>
    <div>
      <h3>Scopes to Request</h3>
      <div v-for="s in oidc.allScopes.value">
        <label :key="s"><input type="checkbox" :value="s" v-model="selectedScopes"> {{ s }}</label>
      </div>
    </div>
    <button @click="loginKeycloak">Begin Login Flow</button>
    or
    <button @click="beginDeviceAuth">Begin Device Flow</button>
    <div v-if="keycloakClaims !== null">
      Logged in as {{ keycloakClaims.name }} ({{ keycloakClaims.email }})
    </div>
  </div>

  <div v-if="beginDeviceResponse !== null">
    <h2>Begin Device Response</h2>
    <pre-json :value="beginDeviceResponse" class="wrap"/>

    <button @click="checkDeviceAuth">Check Device Token</button>
    <div v-if="checkDeviceResponse !== null">
      <h2>Check Device Response</h2>
      <pre-json :value="checkDeviceResponse" class="wrap"/>
    </div>
    <div v-if="checkDeviceError !== null">
      <h2>Check Device Error</h2>
      <pre class="wrap">{{ checkDeviceError }}</pre>
    </div>
  </div>

  <div v-if="accessToken !== null">
    <h2>Access Token</h2>
    <pre class="wrap">{{ accessToken }}</pre>
  </div>
  <div v-if="refreshToken !== null">
    <h2>Refresh Token</h2>
    <pre class="wrap">{{ refreshToken }}</pre>
  </div>

  <div v-if="keycloakClaims !== null">
    <h2>Claims</h2>
    <pre class="wrap">{{ keycloakClaims }}</pre>
  </div>

  <div v-if="accessToken !== null">
    <button @click="logoutKeycloak">Log out</button>
    <button v-if="refreshToken !== null" @click="doRefreshToken">Refresh Token</button>
  </div>
  <div v-if="refreshTokenResponse !== null">
    <h2>Refresh Token Response</h2>
    <pre-json :value="refreshTokenResponse" class="wrap"/>
  </div>
  <div v-if="accessToken !== null">
    <GetTest :token="accessToken"/>
    <UpdateTest :token="accessToken"/>
  </div>
</template>

<style>
.wrap {
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
