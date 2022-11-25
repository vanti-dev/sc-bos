<template>
  <div class="d-flex pa-2">
    <Search />
    <Status />
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

