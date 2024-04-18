<template>
  <component :is="opts.layout" v-bind="opts"/>
</template>

<script setup>
import useDashPage from '@/routes/ops/overview/pages/dashPage.js';
import {reactive} from 'vue';

const props = defineProps({
  pathSegments: {
    type: Array,
    required: true
  }
});

const opts = reactive(useDashPage(() => props.pathSegments));
</script>
