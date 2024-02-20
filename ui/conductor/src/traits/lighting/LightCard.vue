<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Lighting</v-subheader>
      <v-list-item class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">Brightness</v-list-item-title>
        <v-list-item-subtitle class="text-capitalize">{{ brightnessLevelString }}</v-list-item-subtitle>
      </v-list-item>
    </v-list>
    <v-progress-linear
        height="34"
        class="mx-4 my-2"
        :value="brightnessLevelNumber"
        background-color="neutral lighten-1"
        color="accent"/>
    <v-card-actions class="px-4">
      <v-btn small color="neutral lighten-1" :disabled="blockActions" elevation="0" @click="doUpdateBrightness(100)">
        On
      </v-btn>
      <v-btn small color="neutral lighten-1" :disabled="blockActions" elevation="0" @click="doUpdateBrightness(0)">
        Off
      </v-btn>
      <v-spacer/>
      <v-btn
          small
          color="neutral lighten-1"
          elevation="0"
          @click="doUpdateBrightness(brightnessLevelNumber + 1)"
          :disabled="blockActions || brightnessLevelNumber >= 100">
        Up
      </v-btn>
      <v-btn
          small
          color="neutral lighten-1"
          elevation="0"
          @click="doUpdateBrightness(brightnessLevelNumber - 1)"
          :disabled="blockActions || brightnessLevelNumber <= 0">
        Down
      </v-btn>
    </v-card-actions>
    <v-progress-linear color="primary" indeterminate :active="lightLoading || presetsLoading"/>
  </v-card>
</template>

<script setup>
import useAuthSetup from '@/composables/useAuthSetup';
import useLightingTrait from '@/traits/lighting/useLightingTrait.js';

const {blockActions} = useAuthSetup();

const props = defineProps({
  name: {
    type: String,
    default: ''
  }
});

const {
  brightnessLevelString,
  brightnessLevelNumber,
  doUpdateBrightness,
  lightLoading,
  presetsLoading
} = useLightingTrait(() => props.name, false);
</script>

<style scoped>
.v-list-item {
  min-height: auto;
}

.v-progress-linear {
  width: auto;
}
</style>
