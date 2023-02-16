<template>
  <div>ZoneMap</div>
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
