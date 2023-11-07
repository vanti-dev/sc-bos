<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Air Quality</v-subheader>
      <v-list-item v-for="(val, key) of airQualityData" :key="key" class="py-1">
        <v-list-item-title class="text-body-small">{{ key }}</v-list-item-title>
        <v-list-item-subtitle class="font-weight-medium">{{ val }}</v-list-item-subtitle>
      </v-list-item>
    </v-list>
    <v-progress-linear
        v-if="hasScore"
        :value="score"
        height="34"
        class="mx-4 my-2"
        background-color="neutral lighten-1"
        :color="scoreColor"/>
  </v-card>
</template>

<script setup>
import {camelToSentence} from '@/util/string';
import {AirQuality} from '@smart-core-os/sc-api-grpc-web/traits/air_quality_sensor_pb';
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object, // of AirQuality.AsObject
    default: () => ({})
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const score = computed(() => {
  if (props && props.value && props.value.score) {
    return props.value.score;
  }
  return 0;
});

const hasScore = computed(() => {
  return props && props.value && props.value.score;
});

const scoreColor = computed(() => {
  if (score.value < 10) {
    return 'error lighten-1';
  } else if (score.value < 50) {
    return 'warning';
  } else if (score.value < 75) {
    return 'secondary';
  } else {
    return 'success lighten-1';
  }
});

const airQualityData = computed(() => {
  if (props && props.value) {
    const data = {};
    Object.entries(props.value).forEach(([key, value]) => {
      if (value !== undefined) {
        switch (key) {
          case 'carbonDioxideLevel':
            if (value > 0) {
              data['CO2'] = Math.round(value) + ' ppm';
            }
            break;
          case 'volatileOrganicCompounds':
            if (value > 0) {
              data['VOC'] = value.toFixed(3) + ' ppm';
            }
            break;
          case 'airPressure':
            if (value > 0) {
              data['Air Pressure'] = Math.round(value) + ' hPa';
            }
            break;
          case 'infectionRisk':
            if (value > 0) {
              data['Infection Risk'] = Math.round(value) + '%';
            }
            break;
          case 'comfort':
            switch (value) {
              case AirQuality.Comfort.COMFORTABLE:
                data['Comfort'] = 'Comfortable';
                break;
              case AirQuality.Comfort.UNCOMFORTABLE:
                data['Comfort'] = 'Uncomfortable';
                break;
              default:
                // do nothing
            }
            break;
          case 'score':
            if (value > 0) {
              data['Air Quality Score'] = Math.round(value) + '%';
            }
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
