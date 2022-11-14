<template>
  <v-text-field
    v-model="search"
    prepend-inner-icon="mdi-magnify"
    label="Search Spaces"
    single-line
    hide-details
    width="100%"
    outlined
    class="textarea pr-md-4"
  />
</template>

<script>
import { useLightingStore } from "@/stores/operate/lighting.js";
import { storeToRefs } from "pinia";
export default {
  setup() {
    const store = useLightingStore();
    const { search } = storeToRefs(store);

    return { store, search };
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
