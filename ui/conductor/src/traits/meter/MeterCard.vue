<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Meter</v-subheader>
      <v-list-item v-for="item of table" :key="item.label" class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">
          {{ item.label }}
        </v-list-item-title>

        <v-list-item-subtitle class="text-end">
          {{ item.value }} {{ item.unit }}
        </v-list-item-subtitle>
      </v-list-item>
    </v-list>

    <v-progress-linear color="primary" indeterminate :active="loading"/>
  </v-card>
</template>

<script setup>
import {useMeterReading} from '@/traits/meter/meter.js';


const props = defineProps({
  value: {
    type: Object, // of type MeterReading.AsObject
    default: () => {
    }
  },
  info: {
    type: Object, // of type MeterReadingInfo.AsObject
    default: () => null
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const {table} = useMeterReading(() => props.value, () => props.info);
</script>

<style scoped>
.v-list-item {
  min-height: auto;
}

.v-progress-linear {
  width: auto;
}
</style>
