<template>
  <v-card elevation="0" tile>
    <v-card-title class="d-flex text-title-caps-large text-neutral-lighten-3">
      <span>Occupancy Sensor</span>
      <v-spacer/>
      <occupancy-history-card end :name="name"/>
    </v-card-title>
    <v-list tile class="ma-0 pa-0">
      <v-list-item class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">State</v-list-item-title>
        <template #append>
          <v-list-item-subtitle
              :class="[
                stateColor, 'text-capitalize text-body-1 py-1 font-weight-medium text-end']">
            {{ stateStr }}
          </v-list-item-subtitle>
        </template>
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
import OccupancyHistoryCard from '@/traits/occupancy/OccupancyHistoryCard.vue';

const props = defineProps({
  value: {
    type: Object, // of Occupancy.AsObject
    default: () => {
    }
  },
  loading: {
    type: Boolean,
    default: false
  },
  name: {
    type: String,
    required: true
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
