<template>
  <v-card :class="props.class" :style="props.style" class="d-flex flex-column">
    <v-toolbar v-if="!hideToolbar" color="transparent">
      <v-toolbar-title class="text-h4">{{ title }}</v-toolbar-title>
    </v-toolbar>
    <v-card-text class="flex-grow-1 d-flex pt-0">
      <air-quality-history-chart class="flex-grow-1 ma-n2" v-bind="$attrs"/>
    </v-card-text>
  </v-card>
</template>

<script setup>
import AirQualityHistoryChart from '@/dynamic/widgets/environmental/AirQualityHistoryChart.vue';

const props = defineProps({
  title: {
    type: String,
    default: 'Air Quality'
  },
  hideToolbar: {
    type: Boolean,
    default: false
  },
  class: {type: [String, Object, Array], default: undefined},
  style: {type: [String, Object, Array], default: undefined}
});
</script>

<style scoped>

</style>