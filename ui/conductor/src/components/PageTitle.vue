<template>
  <span :class="{hasPageTitle}">{{ pageTitle }}</span>
</template>

<script setup>
import {computed} from 'vue';
import {usePage} from './page.js';

const {pageTitle} = usePage();
const hasPageTitle = computed(() => Boolean(pageTitle));
</script>

<style scoped>
span {
  font-weight: bold;
  font-size: 22px;
}

.hasPageTitle {
  margin-left: 16px;
}
</style>
