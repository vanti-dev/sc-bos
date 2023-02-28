<template>
  <content-card class="mb-8 d-flex flex-column px-6 pt-md-6">
    <h4 class="text-h4 order-lg-last pb-4 pb-lg-0 pt-0 pt-lg-4">Occupancy</h4>
    <v-progress-linear height="24" class="mb-3" v-model="occupancyPercentage"/>
    <div>
      <div class="text-h2" style="float:left">
        {{ occupancy }}
        <span class="text-caption neutral--text text--lighten-5">/{{ maxOccupancy }}</span>
      </div>
      <div class="text-h2" style="float:right">{{ occupancyPercentage.toFixed(0) }}%</div>
    </div>
  </content-card>
</template>

<script setup>

import {computed} from 'vue';
import ContentCard from '@/components/ContentCard.vue';

const props = defineProps({
  occupancy: {
    type: Number,
    default: 0
  },
  maxOccupancy: {
    type: Number,
    default: 1
  }
});

const occupancyPercentage = computed(() => props.occupancy/props.maxOccupancy*100);

</script>

<style scoped>

</style>
