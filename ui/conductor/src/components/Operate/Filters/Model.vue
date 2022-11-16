<template>
  <v-select
    v-model="model"
    :items="models"
    label="Type"
    class="pr-md-4"
    outlined
    style="width: 100px"
  ></v-select>
</template>

<script>
import { useLightingStore } from "@/stores/operate/lighting.js";
import { storeToRefs } from "pinia";
export default {
  setup() {
    const store = useLightingStore();

    const { model, models } = storeToRefs(store);

    return { store, model, models };
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
