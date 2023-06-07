<template>
  <content-card class="mb-8 d-flex flex-column px-6 pt-md-6 pb-md-8 pb-sm-8 ">
    <h4 class="text-h4 py-0 pt-lg-2 pb-4 pl-4">Energy</h4>
    <energy-graph
        class="flex-grow-1 d-none d-md-block"
        width="100%"
        :generated="props.generated"
        :metered="props.metered"/>
    <v-row class="d-flex flex-row justify-center mt-10 mb-1">
      <v-col cols="auto" class="text-h1 align-self-center" style="line-height: 0.3em;">
        <WithElectricDemand
            v-slot="{resource}"
            :name="props.generated">
          {{ storeEnergyValues('generated', resource.value) }}<span style="font-size: 0.5em;">kW</span><br>
        </WithElectricDemand>
        <span class="pl-1 text-title orange--text" style="line-height: 0.35em;">Generated</span>
      </v-col>
      <v-col
          cols="1"
          class="text-h1 d-flex flex-row justify-space-around"
          style="line-height: 0.75em;">
        +
      </v-col>

      <v-col cols="auto" class="text-h1 align-self-center" style="line-height: 0.3em;">
        <WithElectricDemand
            v-slot="{resource}"
            :name="props.metered">
          {{ storeEnergyValues('metered', resource.value) }}<span style="font-size: 0.5em;">kW</span><br>
        </WithElectricDemand>
        <span class="pl-1 text-title primary--text" style="line-height: 0.35em;">Metered</span>
      </v-col>

      <v-col
          cols="1"
          class="text-h1 d-flex flex-row justify-space-around"
          style="line-height: 0.75em;">
        =
      </v-col>

      <v-col cols="auto" class="text-h1 align-self-center" style="line-height: 0.3em;">
        {{ energy.total }}<span style="font-size: 0.5em;">kW</span><br>
        <span class="pl-1 text-title" style="line-height: 0.35em;">Total</span>
      </v-col>
    </v-row>
  </content-card>
</template>

<script setup>
import {reactive, watch} from 'vue';

import ContentCard from '@/components/ContentCard.vue';
import EnergyGraph from '@/routes/ops/components/EnergyGraph.vue';
import WithElectricDemand from '@/routes/devices/components/renderless/WithElectricDemand.vue';

const props = defineProps({
  generated: {
    type: String,
    default: ''
  },
  metered: {
    type: String,
    default: 'building'
  }
});

const energy = reactive({
  generated: 0,
  metered: 0,
  total: 0
});

/**
 *
 * @param {string} type
 * @param {number} value
 * @return {number}
 */
function storeEnergyValues(type, value) {
  if (value) {
    energy[type] = Number((value.realPower / 1000).toFixed(2));
  }

  return energy[type];
};

watch(energy, (newEnergy) => {
  energy.total = (newEnergy.generated + newEnergy.metered).toFixed(2);
}, {immediate: true, deep: true, flush: 'sync'});
</script>

<style scoped>

</style>
