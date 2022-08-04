<template>
  <v-app>
    <v-app-bar app>
      <v-btn>...</v-btn>
      <v-app-bar-title>Smart Core</v-app-bar-title>
      <v-spacer/>
      <v-btn>{{ loginText }}</v-btn>
    </v-app-bar>
    Vue app

    <v-btn @click="count++">Count is {{ count }}</v-btn>
  </v-app>
</template>

<script setup>
import {computed, ref} from 'vue';
import {useAccountStore} from './stores/account.js';

const count = ref(1);

const accountStore = useAccountStore();
const loginText = computed(() => {
  if (accountStore.loggedIn) {
    return accountStore.account.name;
  } else {
    return 'Log in';
  }
})
</script>

<style scoped>
</style>
