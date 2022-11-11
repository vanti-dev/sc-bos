<template>
  <div>Filters</div>
</template>

<script>
import { useOperateStore } from "@/stores/operate.js";
export default {
  setup() {
    const store = useOperateStore();

    return { store };
  },
};
</script>

<style lang="scss" scoped></style>
