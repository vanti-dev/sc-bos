<template>
  <v-card class="d-flex pa-2" elevation="0" height="72px">
    <Search class="mr-4" />
    <Status class="mr-4" />
    <Model />
  </v-card>
</template>

<script setup>
import {useLightingStore} from "@/stores/devices/lighting.js";
import Search from "./Filters/Search.vue";
import Status from "./Filters/Status.vue";
import Model from "./Filters/Model.vue";

const store = useLightingStore();
</script>
