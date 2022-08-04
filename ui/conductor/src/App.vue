<template>
  <v-app>
    <v-app-bar app>
      <v-btn>...</v-btn>
      <v-app-bar-title>Smart Core</v-app-bar-title>
      <v-spacer/>
      <v-btn>Account</v-btn>
    </v-app-bar>
    Vue app

    <v-btn @click="count++">Count is {{ count }}</v-btn>
  </v-app>
</template>

<script setup>
import {ref} from 'vue';

const count = ref(1);
</script>

<style scoped>
</style>
