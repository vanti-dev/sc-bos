<template>
  <div class="text-h1 d-inline-flex justify-space-between" :style="{minWidth: expectedWidth}">
    <labelled-unit
        :value="generatedKW"
        :error="_generated.streamError"
        label="Generated"
        unit="kW"
        label-color="success--text text--lighten-3"/>
    <span v-if="showTotal" class="add mx-3 text-h2">+</span>
    <labelled-unit
        :value="meteredKW"
        :error="_metered.streamError"
        label="Metered"
        unit="kW"
        label-color="primary--text"/>
    <span v-if="showTotal" class="eq mx-3 text-h2">=</span>
    <labelled-unit v-if="showTotal" :value="totalKW" label="Total" unit="kW"/>
    <v-divider v-if="hasOccupancy" vertical class="mx-4"/>
    <labelled-unit
        v-if="hasOccupancy"
        :value="intensityKW"
        :error="_occupancy.streamError"
        label="Energy Intensity"
        unit="kW/person"
        label-color="orange--text"/>
  </div>
</template>

<script setup>
import {usePullElectricDemand} from '@/traits/electricDemand/electric.js';
import {usePullOccupancy} from '@/traits/occupancy/occupancy.js';
import LabelledUnit from '@/widgets/power-total/LabelledUnit.vue';
import useValueOrQuery from '@/widgets/power-total/valueOrQuery.js';
import {computed, reactive} from 'vue';

const props = defineProps({
  generated: {
    type: [
      String, // name of the device
      Object // ElectricDemand.AsObject
    ],
    default: null
  },
  metered: {
    type: [
      String, // name of the device
      Object // ElectricDemand.AsObject
    ],
    default: null
  },
  occupancy: {
    type: [
      String, // name of the device
      Object // Occupancy.AsObject
    ],
    default: null
  }
});

const hasGenerated = computed(() => props.generated != null);
const hasMetered = computed(() => props.metered != null);
const hasOccupancy = computed(() => props.occupancy != null);

const _generated = reactive(useValueOrQuery(() => props.generated, (s) => usePullElectricDemand(s)));

const _metered = reactive(useValueOrQuery(() => props.metered, (s) => usePullElectricDemand(s)));
const _occupancy = reactive(useValueOrQuery(() => props.occupancy, (s) => usePullOccupancy(s)));

const showTotal = computed(() => hasGenerated.value && hasMetered.value);
const totalPower = computed(() => {
  return (_metered.value?.realPower ?? 0) + Math.abs(_generated.value?.realPower ?? 0);
});

const divIfPresent = (a, b) => {
  if (a == null) {
    return a;
  }
  if (b == null || b === 0) {
    return a;
  }
  return a / b;
};
const generatedKW = computed(() => divIfPresent(_generated.value?.realPower, 1000));
const meteredKW = computed(() => divIfPresent(_metered.value?.realPower, 1000));
const totalKW = computed(() => divIfPresent(totalPower.value, 1000));
const intensityKW = computed(() => divIfPresent(totalKW.value, _occupancy.value?.peopleCount));

const segmentWidth = (title, unit = 'kW') => {
  return Math.max(
      title.length * 0.4,
      '0.00'.length + unit.length / 2
  );
};
const expectedWidth = computed(() => {
  let chars = 0;
  if (hasGenerated.value) chars += segmentWidth('Generated');
  if (hasMetered.value) chars += segmentWidth('Metered');
  if (showTotal.value) chars += segmentWidth('Total') + '+'.length + '='.length;
  if (hasOccupancy.value) chars += 1 + segmentWidth('Energy Intensity', 'kW/person');
  return `${chars * 0.6}em`;
});
</script>

<style scoped>

</style>
