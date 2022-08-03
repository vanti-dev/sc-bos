<template>
  <v-app>
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
