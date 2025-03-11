<template>
  <v-card elevation="0" tile>
    <v-card-title class="d-flex text-title-caps-large text-neutral-lighten-3">
      <span>Meter</span>
      <v-spacer/>
      <meter-history-card end :name="name"/>
    </v-card-title>
    <v-list tile class="ma-0 pa-0">
      <v-list-item v-for="item of table" :key="item.label" class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">
          {{ item.label }}
        </v-list-item-title>

        <template #append>
          <v-list-item-subtitle class="text-body-1">
            {{ item.value }} {{ item.unit }}
          </v-list-item-subtitle>
        </template>
      </v-list-item>
    </v-list>

    <v-progress-linear color="primary" indeterminate :active="loading"/>
  </v-card>
</template>

<script setup>
import {useMeterReading} from '@/traits/meter/meter.js';
import MeterHistoryCard from '@/traits/meter/MeterHistoryCard.vue';


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
  },
  name: {
    type: String,
    required: true
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
