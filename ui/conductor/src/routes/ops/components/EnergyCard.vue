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
        {{ currentEnergy.toFixed(1) }}<span style="font-size: 0.5em;">kWh</span><br>
        <span class="pl-1 text-title orange--text" style="line-height: 0.35em;">Generated</span>
      </v-col>
      <v-col
          cols="1"
          class="text-h1 d-flex flex-row align-self-center justify-space-around"
          style="line-height: 0.75em;">
        +
      </v-col>

      <v-col cols="auto" class="text-h1 align-self-center" style="line-height: 0.3em;">
        {{ currentEnergy.toFixed(1) }}<span style="font-size: 0.5em;">kWh</span><br>
        <span class="pl-1 text-title primary--text" style="line-height: 0.35em;">Metered</span>
      </v-col>

      <v-col
          cols="1"
          class="text-h1 d-flex flex-row align-self-center justify-space-around"
          style="line-height: 0.75em;">
        =
      </v-col>

      <v-col cols="auto" class="text-h1 align-self-center" style="line-height: 0.3em;">
        {{ currentEnergy.toFixed(1) }}<span style="font-size: 0.5em;">kWh</span><br>
        <span class="pl-1 text-title" style="line-height: 0.35em;">Total</span>
      </v-col>
    </v-row>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import EnergyGraph from '@/routes/ops/components/EnergyGraph.vue';
import {ref} from 'vue';

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

const currentEnergy = ref(0);
</script>

<style scoped>

</style>
