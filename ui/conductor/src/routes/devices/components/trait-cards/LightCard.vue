<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">State</v-subheader>
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
        color="accent"/>
  </v-card>
</template>

<script setup>

import {computed, defineProps, onUnmounted, reactive, watch} from 'vue';
import {closeResource, newResourceValue} from '@/api/resource';
import {pullBrightness} from '@/api/sc/traits/light';

const props = defineProps({
  name: {
    type: String,
    default: ''
  }
});

const lightValue = reactive(newResourceValue());

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
    return (lightValue.value.levelPercent*100).toFixed(1);
  }
  return '-';
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
