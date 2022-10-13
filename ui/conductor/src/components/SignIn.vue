<template>
  <LocalLogin v-if="loginForm" />
  <LoginChoice v-else />
</template>

<script>
import { useAccountStore } from "@/stores/account.js";
import { storeToRefs } from "pinia";

export default {
  setup() {
    const store = useAccountStore();

    const { loginForm } = storeToRefs(store);

    return { loginForm };
  },
};
</script>

<style lang="scss" scoped></style>
