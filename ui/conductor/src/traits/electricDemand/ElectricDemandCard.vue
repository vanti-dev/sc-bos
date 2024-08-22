<template>
  <v-card elevation="0" tile>
    <v-card-subtitle class="text-title-caps-large text-neutral-lighten-3 py-3 opacity-100">
      Electric
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
import {computed} from 'vue';

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
  apparentPower, apparentPowerUnit,
  reactivePower, reactivePowerUnit,
  powerFactor
} = useElectricDemand(() => props.value);

const rows = computed(() => {
  return [
    {
      label: 'Real Power',
      value: realPower.value.toFixed(3),
      unit: realPowerUnit.value
    },
    {
      label: 'Apparent Power',
      value: apparentPower.value.toFixed(3),
      unit: apparentPowerUnit.value
    },
    {
      label: 'Reactive Power',
      value: reactivePower.value.toFixed(3),
      unit: reactivePowerUnit.value
    },
    {
      label: 'Power Factor',
      value: powerFactor.value?.toFixed(2),
      unit: undefined
    }
  ];
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
