<template>
  <content-card class="mb-5 d-flex flex-column pt-7 pb-0">
    <h4 class="text-h4 pl-4 pb-8 pt-0">Environmental</h4>
    <div :class="['d-flex flex-row justify-center ml-n3', {'flex-wrap mb-4': props.shouldWrap}]">
      <v-col cols="auto" class="ma-0 pa-0">
        <circular-gauge
            v-if="hasTrait(indoorTrait) && getTraitValue(indoorTrait, 'temperatureValue') > 0 || props.shouldWrap"
            :value="getTraitValue(indoorTrait, 'temperatureValue')"
            :color="props.gaugeColor"
            :min="temperatureRange.low"
            :max="temperatureRange.high"
            segments="30"
            style="max-width: 140px;"
            class="mt-2 mb-5 ml-3 mr-2">
          <span class="mt-n4 ml-1 text-h1">
            {{ getTraitValue(indoorTrait, 'temperatureValue').toFixed(1) }}&deg;
          </span>
          <template #title>
            <span class="ml-n1 mb-2">Avg. Indoor Temperature</span>
          </template>
        </circular-gauge>
      </v-col>
      <v-col cols="auto" class="mt-auto mb-0 pb-2 px-0">
        <div
            v-if="hasTrait(outdoorTrait) && getTraitValue(outdoorTrait, 'temperatureValue') > 0 || props.shouldWrap"
            :class="[hasTrait(indoorTrait) && getTraitValue(indoorTrait, 'humidityValue') > 0 ? 'mb-7' : 'mb-2',
                     'd-flex flex-column justify-end align-center']"
            style="width: 150px;">
          <span
              class="text-h1 align-left mb-3"
              style="display: inline-block;">{{ getTraitValue(outdoorTrait, 'temperatureValue').toFixed(1) }}&deg;
          </span>
          <span
              class="text-title text-center"
              style="display: inline-block; width: 100px;">
            External Temperature
          </span>
        </div>
      </v-col>
      <v-col cols="auto" class="pa-0">
        <circular-gauge
            v-if="hasTrait(indoorTrait) && getTraitValue(indoorTrait, 'humidityValue') > 0"
            :value="getTraitValue(indoorTrait, 'humidityValue')"
            :color="props.gaugeColor"
            segments="30"
            style="max-width: 140px;"
            class="mt-2">
          <span class="align-baseline text-h1 mt-n2">
            {{ (getTraitValue(indoorTrait, 'humidityValue') * 100).toFixed(1) }}<span style="font-size: 0.7em;">%</span>
          </span>
          <template #title>
            <span class="mb-2">Avg. Humidity</span>
          </template>
        </circular-gauge>
      </v-col>
    </div>
  </content-card>
</template>

<script setup>
import CircularGauge from '@/components/CircularGauge.vue';
import ContentCard from '@/components/ContentCard.vue';

import {useErrorStore} from '@/components/ui-error/error';
import useAirTemperatureTrait from '@/composables/traits/useAirTemperatureTrait';
import {onBeforeUnmount, ref, watch} from 'vue';

const props = defineProps({
  // name of the device/zone to query for internal temperature data
  name: {
    type: String,
    default: ''
  },
  // name of the device/zone to query for external temperature data
  externalName: {
    type: String,
    default: ''
  },
  gaugeColor: {
    type: String,
    default: ''
  },
  shouldWrap: {
    type: Boolean,
    default: false
  }
});

// todo: do we need to get this from somewhere?
const temperatureRange = ref({
  low: 18.0,
  high: 24.0
});


// Error handling
const errorStore = useErrorStore();
const unwatchErrorFunctions = [];

// Reactive indoor and outdoor traits
const indoorTrait = ref(null);
const outdoorTrait = ref(null);


const getTrait = (name) => {
  if (props[name]) {
    const trait = useAirTemperatureTrait({name: props[name], paused: false});
    unwatchErrorFunctions.push(errorStore.registerValue(trait.airTemperatureResource));
    return trait;
  }
  return null;
};

const hasTrait = (trait) => trait !== null;

const getTraitValue = (trait, key) => hasTrait(trait) ? trait[key] : 0;

const updateTraits = () => {
  if (indoorTrait.value) indoorTrait.value.clearResourceError();
  if (outdoorTrait.value) outdoorTrait.value.clearResourceError();

  unwatchErrorFunctions.forEach(unwatch => unwatch());
  indoorTrait.value = getTrait('name');
  outdoorTrait.value = getTrait('externalName');
};

// Watchers to update traits when props change
watch(() => props.name, updateTraits, {immediate: true});
watch(() => props.externalName, updateTraits, {immediate: true});

// ------------------------------------ //
// Clean up UI Error handling
// Clean up error watchers when the component is unmounted
const cleanup = () => {
  if (indoorTrait.value) indoorTrait.value.clearResourceError();
  if (outdoorTrait.value) outdoorTrait.value.clearResourceError();
  unwatchErrorFunctions.forEach(unwatch => unwatch());
};

onBeforeUnmount(() => {
  cleanup();
});
</script>
