<template>
  <v-main>
    <v-container class="d-flex flex-column align-center my-auto pa-6 fill-height">
      <v-card class="ma-auto pa-4" min-width="450px" max-width="500px">
        <v-card-title class="justify-center text-h1 font-weight-semibold">
          <brand-logo outline="white" style="height: 65px;"/>
          Smart Core
        </v-card-title>

        <local-login v-if="displayLoginForm"/>
        <device-flow-login v-else-if="displayDeviceLogin"/>
        <login-choice v-else @choose="chooseProvider"/>

        <v-card-actions v-if="choiceExists && !displayChoice" class="d-flex justify-center mt-8">
          <a @click="showChoice" class="text-center">Use a different sign in method</a>
        </v-card-actions>
      </v-card>
      <v-btn
          v-if="uiConfig.auth.disableAuthentication"
          class="mt-4"
          color="neutral"
          elevation="0"
          to="/">
        <v-icon class="ml-n2">mdi-chevron-left</v-icon>
        Return to home
      </v-btn>

      <v-snackbar v-model="snackbar.visible">
        {{ snackbar.message }}

        <template #action="{ attrs }">
          <v-btn color="pink" text v-bind="attrs" @click="snackbar.visible = false">
            Close
          </v-btn>
        </template>
      </v-snackbar>
    </v-container>
  </v-main>
</template>
<script setup>
import BrandLogo from '@/components/BrandLogo.vue';
import DeviceFlowLogin from '@/routes/login/DeviceFlowLogin.vue';
import LocalLogin from '@/routes/login/LocalLogin.vue';
import LoginChoice from '@/routes/login/LoginChoice.vue';
import {useAccountStore} from '@/stores/account.js';
import {useUiConfigStore} from '@/stores/ui-config';
import {storeToRefs} from 'pinia';
import {computed, ref} from 'vue';

const uiConfig = useUiConfigStore();
const accountStore = useAccountStore();
const {snackbar} = storeToRefs(accountStore);

const manualDisplayLoginForm = ref(false);
const displayLoginForm = computed(() => {
  return manualDisplayLoginForm.value || accountStore.isOnlyProvider('localAuth');
});

const manualDisplayDeviceLogin = ref(false);
const displayDeviceLogin = computed(() => {
  return manualDisplayDeviceLogin.value || accountStore.isOnlyProvider('deviceFlow');
});

// don't need a display ref for keycloak because it uses redirect, aka there is no page for it

const choiceExists = computed(() => {
  return accountStore.availableAuthProviders.length > 1;
});
const displayChoice = computed(() => {
  return !displayLoginForm.value && !displayDeviceLogin.value;
});

const showChoice = () => {
  manualDisplayLoginForm.value = false;
  manualDisplayDeviceLogin.value = false;
};
const chooseProvider = (p) => {
  switch (p) {
    case 'localAuth':
      manualDisplayLoginForm.value = true;
      break;
    case 'deviceFlow':
      manualDisplayDeviceLogin.value = true;
      break;
    case 'keyCloakAuth':
      // no page to display, redirect instead
      accountStore.loginWithKeyCloak(['profile', 'roles']);
      break;
  }
};

</script>
