<template>
  <div class="d-flex row my-2 mx-1">
    <System />
    <v-spacer />
    <OS />
  </div>
</template>

<script>
import { useLightingStore } from "@/stores/operate/lighting.js";
import System from "./Header/System.vue";
import OS from "./Header/OS.vue";

export default {
  setup() {
    const store = useLightingStore();

    return { store };
  },
};
</script>

<style lang="scss" scoped></style>
