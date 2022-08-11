<template>
  <v-tabs :color="themeColor">
    <v-tab to="/commission/spaces">
      <v-icon left>mdi-domain</v-icon>
      Spaces
    </v-tab>
    <v-tab to="/commission/devices">
      <v-icon left>mdi-cellphone-link</v-icon>
      Devices
    </v-tab>
  </v-tabs>
</template>

<script setup>
import {usePage} from '../../components/page.js';

const {themeColor} = usePage();
</script>

<style scoped>
.v-tabs {
  width: auto;
  flex: initial;
}
</style>
