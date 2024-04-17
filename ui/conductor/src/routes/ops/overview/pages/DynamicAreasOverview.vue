<template>
  <layout-main-side v-bind="opts"/>
</template>

<script setup>
import LayoutMainSide from '@/layout/LayoutMainSide.vue';
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
