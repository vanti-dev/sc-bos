<template>
  <component :is="opts.layout" v-bind="opts"/>
</template>

<script setup>
import useBuildingConfig from '@/routes/ops/overview/pages/buildingConfig.js';
import {reactive} from 'vue';

const opts = reactive(useBuildingConfig());
</script>
