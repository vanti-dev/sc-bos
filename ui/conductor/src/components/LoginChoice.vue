<template>
  <v-card  class="pa-4 ">
    <div class="d-flex justify-end">
      <v-btn @click="store.toggleLoginDialog()" text dense>
        <v-icon> mdi-close </v-icon>
      </v-btn>
    </div>
    <v-card-title class="pt-8 justify-center text-h5 font-weight-semibold">
      Sign in to Smart Core
    </v-card-title>
    <v-card-text class="text-center">
      You can sign in using Keycloak or sign in locally.
    </v-card-text>
    <v-card-actions>
      <v-btn
        @click="store.login(['profile'])"
        color="primary"
        block
        large
        class="font-weight-bold"
        >Sign in with Keycloak</v-btn
      >
    </v-card-actions>
    <v-card-actions>
      <v-btn
        color="warning"
        block
        large
        class="font-weight-bold"
        @click="store.toggleLoginForm()"
        >Use a local Account</v-btn
      >
    </v-card-actions>
    <v-card-text class="text-center">
      Local accounts are used to setup Smart Core, prefer signing in with
      Keycloak.
    </v-card-text>
  </v-card>
</template>

<script setup>
import { useAccountStore } from "@/stores/account.js";

const store = useAccountStore();
</script>

