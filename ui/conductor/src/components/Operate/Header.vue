<template>
  <div class="d-flex row my-2 mx-1">
    <System />
    <v-spacer />
    <OS />
  </div>
</template>

<script>
import { useOperateStore } from "@/stores/operate.js";
import System from "./Header/System.vue";
import OS from "./Header/OS.vue";

export default {
  setup() {
    const store = useOperateStore();

    return { store };
  },
};
</script>

<style lang="scss" scoped></style>
