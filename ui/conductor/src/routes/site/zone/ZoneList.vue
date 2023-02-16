<template>
  <div>{{ zone?.deviceIds ?? 'none selected' }}</div>
</template>

<script setup>
import {Zone} from '@/routes/site/zone/zone';

defineProps({
  zone: {
    type: Zone,
    default: () => {
      return null;
    }
  }
});
</script>

<style scoped>

</style>
