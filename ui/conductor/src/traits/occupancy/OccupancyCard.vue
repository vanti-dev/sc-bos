<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-list-subheader class="text-title-caps-large text-neutral-lighten-3">Occupancy Sensor</v-list-subheader>
      <v-list-item class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">State</v-list-item-title>
        <v-list-item-subtitle
            :class="[
              stateColor, 'text-capitalize text-subtitle-2 py-1 font-weight-medium text-end']">
          {{ stateStr }}
        </v-list-item-subtitle>
      </v-list-item>
      <v-list-item class="py-1" v-if="peopleCount !== 0">
        <v-list-item-title class="text-body-small text-capitalize">Count</v-list-item-title>
        <template #append>
          <v-list-item-subtitle class="text-capitalize text-body-1">{{ peopleCount }}</v-list-item-subtitle>
        </template>
      </v-list-item>
      <v-progress-linear color="primary" indeterminate :active="props.loading"/>
    </v-list>
  </v-card>
</template>

<script setup>
import {useOccupancy} from '@/traits/occupancy/occupancy.js';

const props = defineProps({
  value: {
    type: Object, // of Occupancy.AsObject
    default: () => {
    }
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const {peopleCount, stateStr, stateColor} = useOccupancy(() => props.value);
</script>

<style scoped>
.v-list-item {
  min-height: auto;
}

.v-list-item__subtitle.occupied {
  color: rgb(var(--v-theme-success-lighten-1)) !important;
}

.v-list-item__subtitle.idle {
  color: rgb(var(--v-theme-info)) !important;
}

.v-list-item__subtitle.unoccupied {
  color: rgb(var(--v-theme-warning)) !important;
}
</style>
