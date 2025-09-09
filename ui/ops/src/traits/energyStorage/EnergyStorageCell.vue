<template>
  <status-alert v-if="props.streamError" icon="mdi-battery" :resource="props.streamError"/>

  <span class="text-no-wrap es-cell" v-else-if="percentage !== undefined && !props.streamError">
    <v-tooltip location="bottom">
      <template #activator="{ props: _props }">
        <span v-bind="_props" class="d-flex align-center">
          <span>{{ Math.round(percentage) }}%</span>
          <v-icon
              end
              size="20"
              :icon="batteryIcon"
              :color="batteryColor"/>
        </span>
      </template>
      <div>
        <div>Energy Level: {{ Math.round(percentage) }}%</div>
        <div v-if="flowStatus">Status: {{ flowStatus }}</div>
        <div v-if="hasVoltage">Voltage: {{ voltage.toFixed(1) }}V</div>
        <div v-if="pluggedIn">Plugged In</div>
      </div>
    </v-tooltip>
  </span>
</template>

<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import {useEnergyStorage} from '@/traits/energyStorage/energyStorage.js';
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object, // of EnergyLevel.AsObject
    default: () => {
    }
  },
  loading: {
    type: Boolean,
    default: false
  },
  streamError: {
    type: Object,
    default: null
  }
});

const {percentage, voltage, hasVoltage, flowStatus, pluggedIn, isCharging} = useEnergyStorage(() => props.value);

const batteryIcon = computed(() => {
  if (isCharging.value) {
    return 'mdi-battery-charging';
  }

  const pct = percentage.value;
  if (pct >= 90) return 'mdi-battery';
  if (pct >= 80) return 'mdi-battery-80';
  if (pct >= 60) return 'mdi-battery-60';
  if (pct >= 40) return 'mdi-battery-40';
  if (pct >= 20) return 'mdi-battery-20';
  return 'mdi-battery-alert';
});

const batteryColor = computed(() => {
  const pct = percentage.value;
  if (isCharging.value) return 'success';
  if (pct < 20) return 'error';
  if (pct < 40) return 'warning';
  return 'success';
});
</script>

<style scoped>
.es-cell {
  display: flex;
  align-items: center;
}
</style>
