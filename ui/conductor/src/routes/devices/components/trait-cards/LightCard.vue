<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Lighting</v-subheader>
      <v-list-item class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">Brightness</v-list-item-title>
        <v-list-item-subtitle class="text-capitalize">{{ brightness }}%</v-list-item-subtitle>
      </v-list-item>
    </v-list>
    <v-progress-linear
        height="34"
        class="mx-4 my-2"
        :value="brightness"
        background-color="neutral lighten-1"
        color="accent">
      <template #default>
        <strong class="white--text">{{ statusBarLevels }}</strong>
      </template>
    </v-progress-linear>
    <v-card-actions class="px-4">
      <v-btn small color="neutral lighten-1" :disabled="brightness > 0" elevation="0" @click="updateLight(100)">
        On
      </v-btn>
      <v-btn small color="neutral lighten-1" :disabled="brightness === 0" elevation="0" @click="updateLight(0)">
        Off
      </v-btn>
      <v-spacer/>
      <v-btn
          small
          color="neutral lighten-1"
          elevation="0"
          @click="updateLight(brightness-1)"
          :disabled="brightness <= 0">
        Down
      </v-btn>
      <v-btn
          small
          color="neutral lighten-1"
          elevation="0"
          @click="updateLight(brightness+1)"
          :disabled="brightness >= 100">
        Up
      </v-btn>
    </v-card-actions>
    <v-progress-linear color="primary" indeterminate :active="updateValue.loading"/>
  </v-card>
</template>

<script setup>

import {computed, onMounted, onUnmounted, reactive, watch} from 'vue';
import {closeResource, newActionTracker, newResourceValue} from '@/api/resource';
import {pullBrightness, updateBrightness} from '@/api/sc/traits/light';
import {useErrorStore} from '@/components/ui-error/error';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  }
});

const lightValue = reactive(newResourceValue());
const updateValue = reactive(newActionTracker());

// UI error handling
const errorStore = useErrorStore();
let unwatchLightError; let unwatchUpdateError;
onMounted(() => {
  unwatchLightError = errorStore.registerValue(lightValue);
  unwatchUpdateError = errorStore.registerTracker(updateValue);
});
onUnmounted(() => {
  if (unwatchLightError) unwatchLightError();
  if (unwatchUpdateError) unwatchUpdateError();
});

// if device name changes
watch(() => props.name, async (name) => {
  // close existing stream if present
  closeResource(lightValue);
  // create new stream
  if (name && name !== '') {
    pullBrightness(name, lightValue);
  }
}, {immediate: true});

onUnmounted(() => {
  closeResource(lightValue);
});

const brightness = computed(() => {
  if (lightValue && lightValue.value) {
    return Math.round(lightValue.value.levelPercent);
  }
  return '-';
});

const statusBarLevels = computed(() => {
  let state = 'OFF';
  if (brightness.value > 0 && brightness.value < 100) state = 'DIMMED TO ' + brightness.value + '%';
  else if (brightness.value === 100) state = 'MAX';
  return state;
});


/**
 * @param {number} value
 */
function updateLight(value) {
  /* @type {UpdateBrightnessRequest.AsObject} */
  const req = {
    name: props.name,
    brightness: {
      levelPercent: Math.min(100, Math.round(value))
    }
  };
  updateBrightness(req, updateValue);
}

</script>

<style scoped>
.v-list-item {
  min-height: auto;
}
.v-progress-linear {
  width: auto;
}
</style>
