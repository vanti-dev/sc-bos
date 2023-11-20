<template>
  <sc-logo v-if="isDefault"/>

  <img
      v-else
      :alt="props.theme.appBranding.brandLogo.altText"
      :src="props.theme.appBranding.brandLogo.src">
</template>

<script setup>
import {computed} from 'vue';

import ScLogo from '@/components/ScLogo.vue';

const props = defineProps({
  theme: {
    type: Object,
    default: () => ({})
  }
});

// Check if the default key and it's value exist and is true
//
const isDefault = computed(() => {
  return !props.theme.appBranding?.brandLogo?.src;
});
</script>
