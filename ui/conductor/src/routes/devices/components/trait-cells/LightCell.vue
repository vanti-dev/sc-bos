<template>
  <v-row class="d-flex flex-row flex-nowrap">
    <v-col
        class="text-caption d-flex flex-row justify-center px-1"
        style="min-width: 2.75em; width: 100%;">
      {{ brightnessHotpoint }}
    </v-col>
    <v-col class="px-1">
      <v-icon :color="props.value.brightness > 0 ? 'yellow' : 'white' " size="20">
        {{ lightingIcon }}
      </v-icon>
    </v-col>
  </v-row>
</template>

<script setup>
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object,
    default: () => {}
  }
});

//
//
// Computed

const brightnessHotpoint = computed(() => {
  if (props.value.brightness === 0) {
    return 'Off';
  } else if (props.value.brightness === 100) {
    return 'Max';
  } else if (props.value.brightness > 0 && props.value.brightness < 100) {
    return `${props.value.brightness}%`;
  }

  return '';
});

const lightingIcon = computed(() => {
  if (brightnessHotpoint.value !== '') {
    if (brightnessHotpoint.value !== 'Off') return 'mdi-lightbulb-on';
    else return 'mdi-lightbulb-outline';
  }

  return '';
});
</script>
