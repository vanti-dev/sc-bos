<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Occupancy Sensor</v-subheader>
      <v-list-item class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">State</v-list-item-title>
        <v-list-item-subtitle
            :class="[
              state.toLowerCase(), 'text-capitalize text-subtitle-2 py-1 font-weight-medium text-end']">
          {{ state }}
        </v-list-item-subtitle>
      </v-list-item>
      <v-list-item class="py-1" v-if="count !== 0">
        <v-list-item-title class="text-body-small text-capitalize">Count</v-list-item-title>
        <v-list-item-subtitle class="text-capitalize">{{ count }}</v-list-item-subtitle>
      </v-list-item>
      <v-progress-linear color="primary" indeterminate :active="props.loading"/>
    </v-list>
  </v-card>
</template>

<script setup>

import {occupancyStateToString} from '@/api/sc/traits/occupancy';
import {computed} from 'vue';

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

const count = computed(() => {
  if (props.value) {
    return props.value.peopleCount;
  }
  return 0;
});

const state = computed(() => {
  if (props.value) {
    return occupancyStateToString(props.value.state);
  }
  return '';
});

</script>

<style scoped>
.v-list-item {
  min-height: auto;
}

.v-list-item__subtitle.occupied {
  color: var(--v-success-lighten1) !important;
}

.v-list-item__subtitle.idle {
  color: var(--v-info-base) !important;
}

.v-list-item__subtitle.unoccupied {
  color: var(--v-warning-base) !important;
}
</style>
