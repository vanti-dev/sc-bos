<template>
  <h1 class="pl-1 text-h1">{{ pageTitle }}</h1>
</template>

<script setup>
import {usePage} from '@/components/page.js';

const {pageTitle} = usePage();
</script>

<style scoped>
</style>
