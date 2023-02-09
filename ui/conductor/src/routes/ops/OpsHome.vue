<template>
  <v-container fluid class="pa-0 pt-2 d-flex flex-column" style="min-width: 270px; max-width: 1200px">
    <div class="d-flex flex-column flex-md-row">
      <h3 class="text-h3 pt-2 pb-6 flex-grow-1">Status: Building Overview</h3>
      <content-card class="mr-8 mb-8 mr-sm-0" style="min-width: 250px;">
        <div class="py-1" style="text-align: center">
          <span class="text-title">Smart Core OS: </span>
          <span class="text-title-bold text-uppercase success--text text--lighten-4">Online</span>
        </div>
      </content-card>
    </div>
    <div class="d-flex flex-column flex-lg-row">
      <div class="flex-grow-1 d-flex flex-column mr-lg-8">
        <content-card class="mb-8 d-flex flex-column px-6 pt-md-6">
          <h4 class="text-h4 order-lg-last pb-4 pb-lg-0 pt-0 pt-lg-4">System Monitor</h4>
          <building/>
        </content-card>
        <content-card class="mb-8 d-flex flex-column px-6 pt-md-6">
          <h4 class="text-h4 order-lg-last pb-4 pb-lg-0 pt-0 pt-lg-4">Energy</h4>
        </content-card>
      </div>
      <div class="d-flex flex-column" style="min-width: 250px;">
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
        <content-card class="mb-8 d-flex flex-column px-6 pt-md-6">
          <h4 class="text-h4 order-lg-last pb-4 pb-lg-0 pt-0 pt-lg-4">Environmental</h4>
          <div class="d-flex flex-column flex-md-row flex-lg-column">
            <circular-gauge
                :value="temperature"
                min="15"
                max="35"
                segments="30"
                class="align-self-center mb-6 mb-md-0 mb-lg-8 mr-md-8 mr-lg-0">
              {{ temperature.toFixed(1) }}&deg;
              <template #title>Avg. Indoor Temperature</template>
            </circular-gauge>
            <div class="align-self-center mb-6 mb-md-0 mb-lg-8 mr-md-8 mr-lg-0" style="width: 205px;">
              <span
                  class="text-title"
                  style="display: inline-block; width: 100px;">External Temperature</span>
              <span
                  class="text-h1"
                  style="display: inline-block; float: right;">{{ externalTemperature.toFixed(1) }}&deg;</span>
            </div>
            <circular-gauge
                :value="humidity"
                max="100"
                segments="30"
                class="align-self-center">
              <span class="align-baseline">
                {{ humidity.toFixed(1) }}<span style="font-size: 0.7em;">&percnt;</span>
              </span>
              <template #title>Avg. Humidity</template>
            </circular-gauge>
          </div>
        </content-card>
      </div>
    </div>
    <!-- todo: remove test slider, hook up real backend -->
    <div>
      <span class="text-title-caps">Test</span>
      <v-slider v-model="sliderVal" max="1" step="0.01"/>
    </div>
  </v-container>
</template>
<script setup>
import ContentCard from '@/components/ContentCard.vue';
import CircularGauge from '@/components/CircularGauge.vue';
import {computed, ref} from 'vue';
import Building from '@/clients/ew/Building_EW.vue';

const sliderVal = ref(0);

const temperature = computed(() => 10+(sliderVal.value*25));
const humidity = computed(() => 20+(sliderVal.value*60));
const externalTemperature = computed(() => temperature.value-7);
const occupancy = computed(() => Math.round(sliderVal.value*maxOccupancy.value));
const maxOccupancy = computed(() => 1556);
const occupancyPercentage = computed(() => occupancy.value/maxOccupancy.value*100);

</script>
