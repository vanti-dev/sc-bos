<template>
  <div>
    <v-card-text class="text-center mb-5">
      You can sign in using Keycloak or sign in locally.
    </v-card-text>
    <v-card-actions>
      <v-btn
          @click="store.login(['profile'])"
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
      Local accounts are used to setup Smart Core, prefer signing in with
      Keycloak.
    </v-card-text>
  </div>
</template>

<script setup>
import {useAccountStore} from '@/stores/account.js';

const store = useAccountStore();
</script>
