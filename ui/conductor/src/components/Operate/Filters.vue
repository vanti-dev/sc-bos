<template>
  <div class="d-flex pa-2">
    <Search />
    <Status />
    <Building />
    <Type />
  </div>
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
