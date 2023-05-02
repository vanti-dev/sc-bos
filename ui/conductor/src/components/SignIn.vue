<template>
  <v-card class="pa-4">
    <div class="d-flex justify-end">
      <v-btn @click="store.toggleLoginDialog()" text dense>
        <v-icon> mdi-close </v-icon>
      </v-btn>
    </div>
    <v-card-title class="justify-center text-h3 font-weight-semibold">
      Sign in to Smart Core
    </v-card-title>

    <LocalLogin v-if="loginForm"/>
    <LoginChoice v-else/>
  </v-card>
</template>

<script setup>
import {useAccountStore} from '@/stores/account.js';
import {storeToRefs} from 'pinia';

const store = useAccountStore();

const {loginForm} = storeToRefs(store);
</script>

<style lang="scss" scoped></style>
