<template>
  <LocalLogin v-if="loginForm"/>
  <LoginChoice v-else/>
</template>

<script setup>
import {useAccountStore} from '@/stores/account.js';
import {storeToRefs} from 'pinia';

const store = useAccountStore();

const {loginForm} = storeToRefs(store);
</script>

<style lang="scss" scoped></style>
