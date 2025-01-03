<template>
  <component :is="pageConfig.layout" v-bind="pageConfig.config"/>
</template>

<script setup>
import usePageConfig from '@/routes/ops/overview/pageConfig.js';
import {reactive} from 'vue';

const props = defineProps({
  pathSegments: {
    type: Array,
    required: true
  }
});

const pageConfig = reactive(usePageConfig(() => props.pathSegments));
</script>
