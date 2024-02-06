<template>
  <div>
    <v-card-text :class="[{'mb-8': uiConfig.config.keycloak }, 'text-center mx-auto']" style="max-width: 320px;">
      {{ keycloakMessage.top }}
    </v-card-text>
    <v-card-actions class="d-flex flex-column align-center justify-center">
      <v-btn
          v-if="uiConfig.config.keycloak"
          @click="store.loginWithKeyCloak(['profile', 'roles'])"
          color="primary"
          block
          large
          class="text-body-1 font-weight-bold">
        Sign in
      </v-btn>
    </v-card-actions>
    <v-card-text class="text-body-2 text-center mt-10 mx-auto" style="max-width: 350px;">
      {{ keycloakMessage.bottom }}
    </v-card-text>
    <v-card-actions class="d-flex flex-column align-center justify-center mt-n2">
      <v-btn
          block
          class="text-body-2 ma-0"
          text
          @click="store.loginFormVisible = !store.loginFormVisible">
        Sign in with local Account
      </v-btn>
    </v-card-actions>
  </div>
</template>

<script setup>
import {useAccountStore} from '@/stores/account.js';
import {useUiConfigStore} from '@/stores/ui-config';
import {computed} from 'vue';

const uiConfig = useUiConfigStore();
const store = useAccountStore();

// Tweak the message depending on whether KeyCloak is enabled or not
const keycloakMessage = computed(() => {
  if (uiConfig.config?.keycloak) {
    return {
      top: 'Please sign in to the Smart Core Operator App to unlock all features.',
      // eslint-disable-next-line max-len
      bottom: 'If you are an administrator or need to setup your building, please sign in with a local account.'
    };
  } else {
    return {
      top: 'Please sign in to the Smart Core Operator App with your local account to unlock all features.',
      bottom: 'Local accounts are used to setup Smart Core.'
    };
  }
});
</script>
