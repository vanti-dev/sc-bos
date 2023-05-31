<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Occupancy Sensor</v-subheader>
      <WithOccupancy v-slot="{value}" :name="props.name">
        <v-list-item class="py-1">
          <v-list-item-title class="text-body-small text-capitalize">State</v-list-item-title>
          <v-list-item-subtitle
              :class="[
                value.occupancyState.toLowerCase(),
                'text-capitalize text-subtitle-2 py-1 font-weight-medium'
              ]">
            {{ value.occupancyState }}
          </v-list-item-subtitle>
        </v-list-item>
        <v-list-item class="py-1" v-if="value.occupantCount !== 0">
          <v-list-item-title class="text-body-small text-capitalize">Count</v-list-item-title>
          <v-list-item-subtitle class="text-capitalize">{{ value.occupantCount }}</v-list-item-subtitle>
        </v-list-item>
        <v-progress-linear color="primary" indeterminate :active="value.occupancyValue.loading"/>
      </WithOccupancy>
    </v-list>
  </v-card>
</template>

<script setup>
import WithOccupancy from '../renderless/WithOccupancy.vue';

const props = defineProps({
  name: {
    type: String,
    default: ''
  }
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
