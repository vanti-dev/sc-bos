<template>
  <v-app>
    <router-view/>
  </v-app>
</template>

<script setup>
import {useAccountStore} from '@/stores/account';

useAccountStore(); // force initialisation of the account pinia store so that it has access to the vue context
</script>

<style scoped>
.v-application {
  max-width: 480px;
  min-height: 853px;
  height: 100%;
  background-image: url('/img/background.png');
  background-size: cover;
  background-position: center;
}
</style>
