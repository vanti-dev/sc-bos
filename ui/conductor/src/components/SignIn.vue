<template>
  <v-card class="pa-4">
    <div v-if="config.disableAuthentication" class="d-flex justify-end">
      <v-btn @click="store.toggleLoginDialog()" text dense>
        <v-icon> mdi-close </v-icon>
      </v-btn>
    </div>
    <v-card-title class="justify-center text-h3 font-weight-semibold">
      Sign in to Smart Core
    </v-card-title>

    <LocalLogin v-if="displayLoginForm"/>
    <LoginChoice v-else/>
  </v-card>
</template>

<script setup>
import {computed} from 'vue';
import {useAccountStore} from '@/stores/account.js';
import {useAppConfigStore} from '@/stores/app-config';
import {storeToRefs} from 'pinia';

const {config} = storeToRefs(useAppConfigStore());
const store = useAccountStore();

const {loginForm} = storeToRefs(store);

const displayLoginForm = computed(() => {
  if (config.value?.keycloak) {
    return loginForm.value;
  } else {
    return true;
  }
});
</script>

<style lang="scss" scoped></style>
