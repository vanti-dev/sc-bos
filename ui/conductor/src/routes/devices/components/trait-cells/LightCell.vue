<template>
  <v-row class="d-flex flex-row flex-nowrap">
    <v-col
        class="text-caption d-flex flex-row justify-center px-1"
        style="min-width: 2.75em; width: 100%;">
      {{ brightnessStr }}
    </v-col>
    <v-col class="px-1">
      <v-icon :color="brightness > 0 ? 'yellow' : 'white' " size="20">
        {{ lightingIcon }}
      </v-icon>
    </v-col>
  </v-row>
</template>

<script setup>
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object, // of type Brightness.AsObject
    default: () => {}
  }
});

//
//
// Computed
const brightness = computed(() => props.value?.levelPercent);
const brightnessStr = computed(() => {
  if (brightness.value === 0) {
    return 'Off';
  } else if (brightness.value === 100) {
    return 'Max';
  } else if (brightness.value > 0 && brightness.value < 100) {
    return `${brightness.value}%`;
  }

  return '';
});

const lightingIcon = computed(() => {
  if (brightness.value === 0) return 'mdi-lightbulb-outline';
  if (brightness.value > 0) return 'mdi-lightbulb-on';
  return '';
});
</script>
