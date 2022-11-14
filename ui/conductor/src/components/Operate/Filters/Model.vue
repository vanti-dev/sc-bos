<template>
  <v-select
    v-model="model"
    :items="[
      'All',
      'Philips LED 1245812',
      'Philips LED 1245813',
      'Philips LED 1245814',
    ]"
    label="Type"
    class="pr-md-4"
    outlined
  ></v-select>
</template>

<script>
import { useLightingStore } from "@/stores/operate/lighting.js";
import { storeToRefs } from "pinia";
export default {
  setup() {
    const store = useLightingStore();

    const { model } = storeToRefs(store);

    return { store, model };
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
