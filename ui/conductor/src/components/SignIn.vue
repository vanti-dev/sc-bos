<template>
  <v-card class="pa-4">
    <div v-if="config.disableAuthentication" class="d-flex justify-end">
      <v-btn @click="store.toggleLoginDialog()" text dense>
        <v-icon> mdi-close</v-icon>
      </v-btn>
    </div>
    <v-card-title class="justify-center text-h1 font-weight-semibold">
      <brand-logo outline="white" style="height: 65px;"/>
      Smart Core
    </v-card-title>

    <LocalLogin v-if="displayLoginForm"/>
    <LoginChoice v-else/>
  </v-card>
</template>

<script setup>
import BrandLogo from '@/components/BrandLogo.vue';
import {computed} from 'vue';
import {useAccountStore} from '@/stores/account.js';
import {useAppConfigStore} from '@/stores/app-config';
import {storeToRefs} from 'pinia';

const {config} = storeToRefs(useAppConfigStore());
const store = useAccountStore();

const {loginForm} = storeToRefs(store);

const displayLoginForm = computed(() => {
  // If KeyCloak config available, we can toggle between login variants
  if (config.value?.keycloak) {
    return loginForm.value;

    // If KeyCloak config not available, we can only use local login
  } else {
    return true;
  }
});
</script>

<style lang="scss" scoped></style>
