<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Air Quality</v-subheader>
      <v-list-item v-for="(val, key) of airQualityData" :key="key" class="py-1">
        <v-list-item-title class="text-body-small">{{ key }}</v-list-item-title>
        <v-list-item-subtitle class="text-capitalize font-weight-medium">{{ val }}</v-list-item-subtitle>
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script setup>
import {camelToSentence} from '@/util/string';
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object, // of AirTemperature.AsObject
    default: () => ({})
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const airQualityData = computed(() => {
  if (props && props.value) {
    const data = {};
    Object.entries(props.value).forEach(([key, value]) => {
      if (value !== undefined) {
        switch (key) {
          case 'carbonDioxideLevel':
            data['CO2'] = value.toFixed(2)+'ppm';
            break;
          case 'volatileOrganicCompounds':
            data['VOC'] = value.toFixed(2)+'ppm';
            break;
          case 'airPressure':
            data['Air Pressure'] = value.toFixed(2)+' hPa';
            break;
          case 'infectionRisk':
            data['Infection Risk'] = Math.round(value*100)+'%';
            break;
          default: {
            data[camelToSentence(key)] = value;
          }
        }
      }
    });
    return data;
  }
  return {};
});

</script>

<style scoped>
.v-list-item {
  min-height: auto;
}

.v-progress-linear {
  width: auto;
}
</style>
