<template>
  <v-card>
    <v-card-text> System: Lighting </v-card-text>
  </v-card>
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
