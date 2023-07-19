<template>
  <v-card elevation="0" tile>
    <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Electric</v-subheader>
    <div class="layout mx-4">
      <template v-for="(row, i) in rows">
        <span :key="i+'label'" class="label text-body-small">{{ row.label }}</span>
        <span :key="i+'value'" class="value">{{ row.value ?? '' }}</span>
        <span :key="i+'unit'" class="unit">{{ row.unit ?? '' }}</span>
      </template>
    </div>
    <v-progress-linear color="primary" indeterminate :active="props.loading"/>
  </v-card>
</template>
<script setup>
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

const realPower = computed(() => props.value?.realPower);
const apparentPower = computed(() => props.value?.apparentPower);
const reactivePower = computed(() => props.value?.reactivePower);
const powerFactor = computed(() => props.value?.powerFactor);

const rows = computed(() => {
  return [
    {
      label: 'Real Power',
      value: realPower.value?.toFixed(3),
      unit: 'kW'
    },
    {
      label: 'Apparent Power',
      value: apparentPower.value?.toFixed(3),
      unit: 'kVA'
    },
    {
      label: 'Reactive Power',
      value: reactivePower.value?.toFixed(3),
      unit: 'kVAr'
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
