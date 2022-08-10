<template>
  <span :class="{hasPageName}">{{ pageName }}</span>
</template>

<script setup>
import {computed} from 'vue';
import {useRoute} from '../util/router.js';

const route = useRoute();
const pageName = computed(() => route?.meta?.title);
const hasPageName = computed(() => Boolean(pageName.value));
</script>

<style scoped>
span {
  font-weight: bold;
  font-size: 22px;
}

.hasPageName {
  margin-left: 16px;
}
</style>
