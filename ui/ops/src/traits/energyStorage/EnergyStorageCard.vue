<template>
  <v-card elevation="0" tile>
    <v-card-subtitle class="text-title-caps-large text-neutral-lighten-3 py-3 opacity-100 d-flex align-center">
      <span class="mr-auto">Energy Storage</span>
      <v-btn v-if="hasHidden"
             :icon="showHidden ? 'mdi-chevron-up' : 'mdi-chevron-down'"
             variant="flat" size="small"
             class="my-n4"
             @click="showHidden = !showHidden"
             v-tooltip="showHidden ? 'Hide more fields' : 'Show more fields'"/>
    </v-card-subtitle>
    <div class="layout mx-4">
      <template v-for="(row, i) in rows" :key="i">
        <span class="label text-body-small">{{ row.label }}</span>
        <span class="value">{{ row.value ?? '' }}</span>
        <span class="unit">{{ row.unit ?? '' }}</span>
      </template>
    </div>
    <v-progress-linear color="primary" indeterminate :active="props.loading"/>
  </v-card>
</template>

<script setup>
import {useEnergyStorage} from '@/traits/energyStorage/energyStorage.js';
import {computed, ref} from 'vue';

const props = defineProps({
  value: {
    type: Object, // of EnergyLevel.AsObject
    default: () => {
    }
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const {
  percentage, energyKwh, distanceKm, voltage, hasVoltage, descriptive,
  flowStatus, speed, targetPercentage, hasTargetPercentage, pluggedIn
} = useEnergyStorage(() => props.value);

const showHidden = ref(false);
const hasHidden = computed(() => !energyKwh.value || !distanceKm.value || !hasVoltage.value || !speed.value);
const showEnergyKwh = computed(() => energyKwh.value > 0 || showHidden.value);
const showDistanceKm = computed(() => distanceKm.value > 0 || showHidden.value);
const showVoltage = computed(() => hasVoltage.value || showHidden.value);
const showSpeed = computed(() => speed.value || showHidden.value);

const rows = computed(() => {
  const result = [
    {
      label: 'Battery Level',
      value: percentage.value.toFixed(1),
      unit: '%'
    },
    {
      label: 'Status',
      value: flowStatus.value,
      unit: undefined
    }
  ];

  if (descriptive.value) {
    result.push({
      label: 'Level',
      value: descriptive.value,
      unit: undefined
    });
  }

  if (showEnergyKwh.value) {
    result.push({
      label: 'Energy',
      value: energyKwh.value.toFixed(1),
      unit: 'kWh'
    });
  }

  if (showVoltage.value) {
    result.push({
      label: 'Voltage',
      value: voltage.value.toFixed(1),
      unit: 'V'
    });
  }

  if (showDistanceKm.value) {
    result.push({
      label: 'Range',
      value: distanceKm.value.toFixed(0),
      unit: 'km'
    });
  }

  if (showSpeed.value) {
    result.push({
      label: 'Speed',
      value: speed.value,
      unit: undefined
    });
  }

  if (hasTargetPercentage.value) {
    result.push({
      label: 'Target',
      value: targetPercentage.value.toFixed(1),
      unit: '%'
    });
  }

  result.push({
    label: 'Plugged In',
    value: pluggedIn.value ? 'Yes' : 'No',
    unit: undefined
  });

  return result;
});
</script>

<style scoped>
.layout {
  display: grid;
  grid-template-columns: auto 1fr auto;
  grid-gap: 4px 0.3em;
  align-items: baseline;
}

.value {
  text-align: right;
}

.value, .unit {
  color: rgba(255, 255, 255, 0.7);
}
</style>
