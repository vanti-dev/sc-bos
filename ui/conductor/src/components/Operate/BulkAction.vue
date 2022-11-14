<template>
  <v-card
    class="pa-2 d-flex"
    flat
    z-index="10"
    color="#111721"
    v-if="selected.length > 1"
  >
    <v-card-subtitle> Bulk Actions</v-card-subtitle>
    <v-card-actions>
      <v-btn color="green" class="mr-2">Turn On</v-btn>
      <v-btn color="red" class="mr-2">Turn Off</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import { useLightingStore } from "@/stores/operate/lighting.js";
import { storeToRefs } from "pinia";

export default {
  setup() {
    const store = useLightingStore();

    const { selected } = storeToRefs(store);

    return { selected };
  },
};
</script>

<style lang="scss" scoped></style>
