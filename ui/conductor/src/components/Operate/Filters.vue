<template>
  <v-layout row wrap ma-4 pt-4>
    <v-flex xs12 sm12 md6 lg6 xl6><Search /> </v-flex>
    <v-flex xs12 sm12 md2 lg2 xl2><Status /> </v-flex>
    <v-flex xs12 sm12 md2 lg2 xl2><Building /></v-flex>
    <v-flex xs12 sm12 md2 lg2 xl2> <Type /></v-flex>
  </v-layout>
</template>

<script>
import { useLightingStore } from "@/stores/operate/lighting.js";
import Search from "./Filters/Search.vue";
import Type from "./Filters/Type.vue";
import Status from "./Filters/Status.vue";
import Building from "./Filters/Building.vue";

export default {
  setup() {
    const store = useLightingStore();

    return { store };
  },
};
</script>

<style lang="scss" scoped></style>
