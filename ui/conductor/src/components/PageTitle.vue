<template>
  <span :class="{hasPageTitle}">{{ pageTitle }}</span>
</template>

<script setup>
import {usePage} from '@/components/page.js';
import {computed} from 'vue';

const {pageTitle} = usePage();
const hasPageTitle = computed(() => Boolean(pageTitle));
</script>

<style scoped>
span {
  font-weight: bold;
  font-size: 14px;
  text-transform: uppercase;
}

.hasPageTitle {
  margin-left: 16px;
}
</style>
