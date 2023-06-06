<template>
  <content-card class="mb-8 d-flex flex-column px-6 pt-md-6 pb-md-8 pb-sm-8 ">
    <h4 class="text-h4 py-0 pt-lg-2 pb-4">Energy</h4>
    <energy-graph class="flex-grow-1 d-none d-md-block" width="100%" :name="props.zoneName"/>
    <div class="d-flex flex-row justify-center">
      <div class="text-h2 align-self-center" style="line-height: 0.5em; min-width: 120px;">
        {{ currentEnergy.toFixed(1) }}<span style="font-size: 0.5em;">kWh</span><br>
        <span class="text-title">Consumption</span>
      </div>
    </div>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import EnergyGraph from '@/routes/ops/components/EnergyGraph.vue';
import {ref} from 'vue';

const props = defineProps({
  zoneName: {
    type: [Array, String],
    default: 'building'
  }
});

const currentEnergy = ref(0);
</script>

<style scoped>

</style>
