<template>
  <v-tooltip location="bottom">
    <template #activator="{ props }">
      <div v-bind="props">
        <v-menu
            location="bottom left"
            :close-on-content-click="false"
            min-width="275px">
          <template #activator="{ props }">
            <v-btn
                theme="dark"
                rounded="circle"
                v-bind="props">
              <v-icon>mdi-dots-vertical</v-icon>
            </v-btn>
          </template>

          <v-list height="100%">
            <v-list-item
                class="d-flex flex-row justify-center pa-0 mt-n1 px-3"
                @click.stop="showConversionToggle = !showConversionToggle">
              <v-list-subheader class="text-body-2 pa-0">Unit Type</v-list-subheader>
              <v-spacer/>
              <v-switch
                  class="ml-4 my-auto no-pointer-events"
                  color="primary"
                  density="compact"
                  readonly
                  hide-details
                  inset
                  :value="props.showConversion">
                <template #prepend>
                  <span class="text-caption text-grey-lighten-1">kW</span>
                </template>
                <template #append>
                  <span class="text-caption text-grey-lighten-1 ml-n4">CO₂</span>
                </template>
              </v-switch>
            </v-list-item>
            <v-list-item class="pa-0 d-flex flex-row justify-center px-3" density="compact">
              <v-list-subheader class="text-body-2 pa-0">Duration</v-list-subheader>
              <v-spacer/>
              <v-btn-toggle
                  v-model="activeDuration"
                  selected-class="primary"
                  density="compact"
                  mandatory>
                <v-btn
                    v-for="option in durationOptions"
                    selected-class="primary text--darken-3"
                    class="bg-transparent text-grey-lighten-1"
                    :key="option.id"
                    size="small"
                    :value="option.value">
                  <span class="text-caption">{{ option.text }}</span>
                </v-btn>
              </v-btn-toggle>
            </v-list-item>
            <v-list-item
                class="pa-0 d-flex flex-row align-left align-center px-3 mb-n1"
                density="compact"
                @click="emits('exportCSV')">
              <v-list-subheader class="text-body-2 pa-0">Export CSV...</v-list-subheader>
            </v-list-item>
          </v-list>
        </v-menu>
      </div>
    </template>
    <span>Options</span>
  </v-tooltip>
</template>

<script setup>
import {DAY, HOUR, MINUTE} from '@/components/now.js';

const emits = defineEmits(['exportCSV']);

// Defining the options for the duration type buttons
const durationOptions = [
  {
    text: '24H',
    value: {
      id: '24H',
      span: 20 * MINUTE,
      timeFrame: 24 * HOUR
    }
  },
  {
    text: '1W',
    value: {
      id: '1W',
      span: 2 * HOUR,
      timeFrame: 7 * DAY
    }
  },
  {
    text: '30D',
    value: {
      id: '30D',
      span: 6 * HOUR,
      timeFrame: 30 * DAY
    }
  }
];

// Computed property to toggle between kW and CO2
// Syncs with the parent component
const showConversionToggle = defineModel('showConversion', {
  type: Boolean,
  default: false
});

// Computed property to toggle between the duration options
// Syncs with the parent component
const activeDuration = defineModel('durationOption', {
  type: Object,
  default: () => durationOptions[0].value
});
</script>

<style lang="scss" scoped>
.no-pointer-events {
  pointer-events: none;
}
</style>
