<script setup>
import {allScopes, keycloak} from "./auth"
import {ref} from "vue"
import GetTest from "./component/GetTest.vue"
import UpdateTest from "./component/UpdateTest.vue";

const selectedScopes = ref([])
const accessToken = ref(null)

const keycloakClaims = ref(null)

keycloak.init({})
    .then(authenticated => {
      if (authenticated) {
        keycloakClaims.value = keycloak.idTokenParsed
        accessToken.value = keycloak.token
      }
    })
    .catch(console.error)

async function loginKeycloak() {
  const scopes = selectedScopes.value.join(" ")

  try {
    await keycloak.login({
      scope: scopes
    })
  } catch (e) {
    console.error(e)
  }
}
</script>

<template>
  <h1>Auth Demo</h1>
  <div>
    <h2>Keycloak</h2>
    <div>
      <h3>Scopes to Request</h3>
      <div v-for="scope in allScopes">
        <input type="checkbox" :id="scope" :value="scope" v-model="selectedScopes">
        <label :for="scope">{{ scope }}</label>
      </div>
    </div>
    <button @click="loginKeycloak">Login</button>
    <div v-if="keycloakClaims !== null">
      Logged in as {{ keycloakClaims.name }} ({{ keycloakClaims.email }})
    </div>
  </div>

  <div v-if="accessToken !== null">
    <h2>Access Token</h2>
    <pre class="wrap">{{ accessToken }}</pre>
  </div>

  <div v-if="accessToken !== null">
    <GetTest :token="accessToken"/>
    <UpdateTest :token="accessToken"/>
  </div>
</template>

<style>
  pre.wrap {
    white-space: pre-wrap;
    word-break: break-all;
  }
</style>
