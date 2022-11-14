<template>
  <v-select :items="['All']" label="Status" class="pr-md-4" outlined></v-select>
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

<style lang="scss" scoped>
.v-text-field--outlined ::v-deep(.v-input__control > .v-input__slot) {
  background-color: #3f454a;
  color: white;
  width: 100%;
}
</style>
