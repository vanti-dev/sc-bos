<template>
  <div class="d-flex pa-2">
    <Search />
    <Status />
    <Model />
  </div>
</template>

<script>
import { useLightingStore } from "@/stores/operate/lighting.js";

export default {
  setup() {
    const store = useLightingStore();

    return { store };
  },
};
</script>

<style lang="scss" scoped></style>
