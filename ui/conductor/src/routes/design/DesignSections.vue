<template>
  <v-tabs :color="themeColor">
    <v-tab to="/design/buildings">
      <v-icon left>mdi-domain</v-icon>
      Buildings & Floors
    </v-tab>
    <v-tab to="/design/spaces">
      <v-icon left>mdi-map-marker</v-icon>
      Spaces
    </v-tab>
    <v-tab to="/design/plan">
      <v-icon left>mdi-view-grid-plus-outline</v-icon>
      Floor Plan Designer
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
