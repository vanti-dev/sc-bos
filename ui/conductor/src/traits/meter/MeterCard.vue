<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Meter</v-subheader>
      <v-list-item class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">
          {{ meterDetails.key }}
        </v-list-item-title>

        <v-list-item-subtitle class="text-end">
          {{ meterDetails.value }} {{ meterDetails.unit }}
        </v-list-item-subtitle>
      </v-list-item>
    </v-list>

    <v-progress-linear color="primary" indeterminate :active="loading"/>
  </v-card>
</template>

<script setup>
import {computed} from 'vue';


const props = defineProps({
  value: {
    type: Object, // of type MeterReading.AsObject
    default: () => {
    }
  },
  unit: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const meterDetails = computed(() => {
  return {
    key: 'Usage',
    value: props.value?.usage?.toFixed(2) ?? '-',
    unit: props.unit
  };
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
