<template>
  <v-card elevation="0" tile>
    <v-card-title class="d-flex text-title-caps-large text-neutral-lighten-3">
      <span>Air Quality</span>
      <v-spacer/>
      <air-quality-history-card end :name="name"/>
    </v-card-title>
    <v-list tile class="ma-0 pa-0">
      <v-list-item v-for="(val, key) of presentMetrics" :key="key" class="py-1">
        <v-list-item-title class="text-body-small">{{ key }}</v-list-item-title>
        <template #append>
          <v-list-item-subtitle class="font-weight-medium text-body-1">{{ val.value }}</v-list-item-subtitle>
        </template>
      </v-list-item>
    </v-list>
    <v-progress-linear
        v-if="score"
        :model-value="score.value"
        height="34"
        class="mx-4 my-2"
        bg-color="neutral-lighten-1"
        bg-opacity="1"
        :color="scoreColor"/>
  </v-card>
</template>

<script setup>
import {useAirQuality, useStatusColor} from '@/traits/airQuality/airQuality.js';
import AirQualityHistoryCard from '@/traits/airQuality/AirQualityHistoryCard.vue';

const props = defineProps({
  value: {
    type: Object, // of AirQuality.AsObject
    default: () => ({})
  },
  loading: {
    type: Boolean,
    default: false
  },
  name: {
    type: String,
    required: true
  }
});

const {score, presentMetrics} = useAirQuality(() => props.value);
const scoreColor = useStatusColor(score);

</script>

<style scoped>
.v-list-item {
  min-height: auto;
}

.v-progress-linear {
  width: auto;
}

</style>
