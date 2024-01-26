<template>
  <div class="d-flex flex-column align-center my-auto">
    <v-card class="ma-auto pa-4" min-width="450px" max-width="500px">
      <v-card-title class="justify-center text-h1 font-weight-semibold">
        <brand-logo outline="white" style="height: 65px;"/>
        Smart Core
      </v-card-title>

      <LocalLogin v-if="displayLoginForm"/>
      <LoginChoice v-else/>
    </v-card>
    <v-btn
        v-if="appConfigStore.config.disableAuthentication"
        class="mt-4"
        color="neutral"
        elevation="0"
        to="/">
      <v-icon class="ml-n2">mdi-chevron-left</v-icon>
      Return to home
    </v-btn>
  </div>
</template>
<script setup>
import BrandLogo from '@/components/BrandLogo.vue';
import LocalLogin from '@/routes/login/LocalLogin.vue';
import LoginChoice from '@/routes/login/LoginChoice.vue';
import {useAccountStore} from '@/stores/account.js';
import {useAppConfigStore} from '@/stores/app-config';
import {computed} from 'vue';

const appConfigStore = useAppConfigStore();
const accountStore = useAccountStore();

const displayLoginForm = computed(() => {
  // If KeyCloak config available, we can toggle between login variants
  if (appConfigStore.config.keycloak) {
    return accountStore.loginFormVisible;

    // If KeyCloak config not available, we can only use local login
  } else {
    return true;
  }
});
</script>
