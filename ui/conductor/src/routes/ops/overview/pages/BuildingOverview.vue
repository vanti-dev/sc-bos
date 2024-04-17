<template>
  <layout-main-side v-bind="opts"/>
</template>

<script setup>
import LayoutMainSide from '@/layout/LayoutMainSide.vue';
import useBuildingConfig from '@/routes/ops/overview/pages/buildingConfig.js';
import {reactive} from 'vue';

const opts = reactive(useBuildingConfig());
</script>

<style scoped>
</style>
