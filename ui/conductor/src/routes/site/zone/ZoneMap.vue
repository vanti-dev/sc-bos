<template>
  <v-container>
    <level0/>
  </v-container>
</template>

<script setup>
import Level0 from '@/clients/ew/Level0.vue';
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
