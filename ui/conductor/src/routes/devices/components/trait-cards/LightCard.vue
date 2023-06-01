<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Lighting</v-subheader>
      <v-list-item class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">Brightness</v-list-item-title>
        <v-list-item-subtitle class="text-capitalize">{{ brightnessStr }}</v-list-item-subtitle>
      </v-list-item>
    </v-list>
    <v-progress-linear
        height="34"
        class="mx-4 my-2"
        :value="brightness"
        background-color="neutral lighten-1"
        color="accent"/>
    <v-card-actions class="px-4">
      <v-btn small color="neutral lighten-1" elevation="0" @click="updateLight(100)">On</v-btn>
      <v-btn small color="neutral lighten-1" elevation="0" @click="updateLight(0)">Off</v-btn>
      <v-spacer/>
      <v-btn
          small
          color="neutral lighten-1"
          elevation="0"
          @click="updateLight(brightness+1)"
          :disabled="brightness >= 100">
        Up
      </v-btn>
      <v-btn
          small
          color="neutral lighten-1"
          elevation="0"
          @click="updateLight(brightness-1)"
          :disabled="brightness <= 0">
        Down
      </v-btn>
    </v-card-actions>
    <v-progress-linear color="primary" indeterminate :active="loading"/>
  </v-card>
</template>

<script setup>

import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object, // of type Brightness.AsObject
    default: () => {
    }
  },
  loading: {
    type: Boolean,
    default: false
  }
});
const emit = defineEmits(['updateBrightness']);
const brightness = computed(() => props.value?.levelPercent ?? 0);
const brightnessStr = computed(() => {
  const val = brightness.value;
  if (val === 0) {
    return 'Off';
  } else if (val === 100) {
    return 'Max';
  } else if (val > 0 && val < 100) {
    return `${val.toFixed(0)}%`;
  }
  return '';
});

/**
 * @param {number} brightness
 */
function updateLight(brightness) {
  emit('updateBrightness', {
    brightness: {
      levelPercent: Math.min(100, Math.round(brightness))
    }
  });
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
