<template>
  <div>
    <v-card-text :class="[{'mb-5': appConfig.config?.keycloak }, 'text-center']">
      {{ keycloakMessage.top }}
    </v-card-text>
    <v-card-actions>
      <v-btn
          v-if="appConfig.config?.keycloak"
          @click="doLogin()"
          color="primary"
          block
          large
          class="text-body-1 font-weight-bold">
        Sign in with Keycloak
      </v-btn>
    </v-card-actions>
    <v-card-actions>
      <v-btn
          color="warning"
          block
          large
          class="text-body-1 font-weight-bold"
          @click="store.toggleLoginForm()">
        Use a local Account
      </v-btn>
    </v-card-actions>
    <v-card-text class="text-center mt-4">
      {{ keycloakMessage.bottom }}
    </v-card-text>
  </div>
</template>

<script setup>
import {computed} from 'vue';
import {useAccountStore} from '@/stores/account.js';
import {useAppConfigStore} from '@/stores/app-config';

const appConfig = useAppConfigStore();
const store = useAccountStore();
const doLogin = () => store.login(['profile', 'roles'])
    .catch(err => console.error(err));

// Tweak the message depending on whether KeyCloak is enabled or not
const keycloakMessage = computed(() => {
  if (appConfig.config?.keycloak) {
    return {
      top: 'You can sign in using Keycloak or sign in locally.',
      bottom: 'Local accounts are used to setup Smart Core, prefer signing in with Keycloak.'
    };
  } else {
    return {
      top: 'You can sign in using a local account.',
      bottom: 'Local accounts are used to setup Smart Core.'
    };
  }
});
</script>
