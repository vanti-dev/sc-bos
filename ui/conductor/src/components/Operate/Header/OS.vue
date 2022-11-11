<template>
  <v-card>
    <v-card-text> Smart Core OS: ONLINE </v-card-text>
  </v-card>
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
