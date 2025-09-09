<template>
  <v-card elevation="0" tile>
    <v-card-subtitle class="text-title-caps-large text-neutral-lighten-3 py-3 opacity-100 d-flex align-center">
      <span class="mr-auto">Electric</span>
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
import {useElectricDemand} from '@/traits/electricDemand/electric.js';
import {computed, ref} from 'vue';

const props = defineProps({
  value: {
    type: Object, // of ElectricDemand.AsObject
    default: () => {
    }
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const {
  realPower, realPowerUnit,
  apparentPower, apparentPowerUnit, hasApparentPower,
  reactivePower, reactivePowerUnit, hasReactivePower,
  powerFactor, hasPowerFactor,
  current, currentUnit, hasCurrent
} = useElectricDemand(() => props.value);

const showHidden = ref(false);
const hasHidden = computed(() => !hasApparentPower.value || !hasReactivePower.value || !hasPowerFactor.value);
const showApparentPower = computed(() => hasApparentPower.value || showHidden.value);
const showReactivePower = computed(() => hasReactivePower.value || showHidden.value);
const showPowerFactor = computed(() => hasPowerFactor.value || showHidden.value);

const rows = computed(() => {
  const result = [
    {
      label: 'Real Power',
      value: realPower.value.toFixed(3),
      unit: realPowerUnit.value
    }
  ];
  if (showApparentPower.value) {
    result.push({
      label: 'Apparent Power',
      value: apparentPower.value.toFixed(3),
      unit: apparentPowerUnit.value
    });
  }
  if (showReactivePower.value) {
    result.push({
      label: 'Reactive Power',
      value: reactivePower.value.toFixed(3),
      unit: reactivePowerUnit.value
    });
  }
  if (showPowerFactor.value) {
    result.push({
      label: 'Power Factor',
      value: powerFactor.value?.toFixed(2),
      unit: undefined
    });
  }
  if (hasCurrent.value) {
    result.push({
      label: 'Current',
      value: current.value.toFixed(3),
      unit: currentUnit.value
    });
  }
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
