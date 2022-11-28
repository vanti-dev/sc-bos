<template>
  <div class="d-flex pa-2">
    <Search class="mr-4" />
    <Status class="mr-4" />
    <Model />
  </div>
</template>

<script setup>
import { useLightingStore } from "@/stores/operate/lighting.js";
import Search from "./Filters/Search.vue";
import Status from "./Filters/Status.vue";
import Model from "./Filters/Model.vue";

const store = useLightingStore();
</script>

