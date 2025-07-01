<template>
  <v-list-item title="Range">
    <template #append>
      <v-btn-toggle
          v-model="periodSelected"
          variant="outlined"
          density="compact"
          divided
          mandatory
          class="ml-4">
        <v-btn
            v-for="item in periodSelections"
            :key="item.text ?? item.icon"
            :value="item"
            :text="item.text"
            :icon="item.icon"
            size="small"
            width="2.5rem"
            min-width="auto"
            class="px-0"
            v-tooltip:bottom="item.tooltip"/>
      </v-btn-toggle>
    </template>
  </v-list-item>
  <v-list-item v-if="!periodSelectionIsDates" title="Offset">
    <template #append>
      <v-btn-toggle
          mandatory v-model="offset"
          variant="outlined" density="compact"
          divided class="ml-4">
        <v-btn :value="-1" size="small" text="Last" v-tooltip:bottom="'Show the previous period'"/>
        <v-btn :value="0" size="small" text="Current Period" v-tooltip:bottom="'Show the current period'"/>
      </v-btn-toggle>
    </template>
  </v-list-item>
  <v-list-item v-if="periodSelectionIsDates">
    <v-date-input
        v-model="periodSelectedDates" multiple="range"
        label="Custom date range" placeholder="from - to" persistent-placeholder
        hide-details/>
  </v-list-item>
</template>

<script setup>
import {isNamedPeriod} from '@/composables/time.js';
import {computed, ref, watch} from 'vue';
import {VDateInput} from 'vuetify/labs/components';

const start = defineModel('start', {type: [Date, String], required: true});
const end = defineModel('end', {type: [Date, String], required: true});
const offset = defineModel('offset', {type: Number, default: 0});

// user editable date selection
const propPeriod = computed(() => {
  if (isNamedPeriod(start.value) && start.value === end.value) {
    return start.value;
  }
  return undefined;
})
/**
 * @typedef {Object} PeriodSelection
 * @property {string} [text] - the text to display on the button
 * @property {string} [value] - the value to set when selected
 * @property {string} [icon] - the icon to display on the button
 * @property {string} [tooltip] - the tooltip to display when hovering over the button
 */

/** @type {PeriodSelection[]} */
const periodSelections = [
  {text: 'D', value: 'day', tooltip: 'Day'},
  {text: 'W', value: 'week', tooltip: 'Week'},
  {text: 'M', value: 'month', tooltip: 'Month'},
  {text: 'Y', value: 'year', tooltip: 'Year'},
  {icon: 'mdi-calendar', tooltip: 'Custom date range'},
];
const periodSelected = ref(periodSelections[0]);
watch(propPeriod, (period) => {
  if (period) {
    const selection = periodSelections.find((s) => s.value === period);
    if (selection) {
      periodSelected.value = selection;
      return;
    }
  }
  periodSelected.value = periodSelections[periodSelections.length - 1];
}, {immediate: true});
const periodSelectionIsDates = computed(() => !Object.hasOwn(periodSelected.value, 'value'));
const periodSelectedDates = ref([]);
watch([periodSelected, periodSelectedDates], ([period, dates]) => {
  if (period.value) {
    start.value = period.value;
    end.value = period.value;
    return;
  }
  if (dates.length === 0) return;
  start.value = dates[0];
  end.value = dates[dates.length - 1];
})
</script>

<style scoped>

</style>