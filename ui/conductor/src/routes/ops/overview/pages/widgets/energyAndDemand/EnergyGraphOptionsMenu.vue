<template>
  <v-tooltip bottom>
    <template #activator="{ on: onTooltip, attrs: attrsTooltip }">
      <div v-bind="attrsTooltip" v-on="onTooltip">
        <v-menu
            bottom
            :close-on-content-click="false"
            left
            min-width="275px"
            offset-y>
          <template #activator="{ on, attrs }">
            <v-btn
                dark
                icon
                v-bind="attrs"
                v-on="on">
              <v-icon>mdi-dots-vertical</v-icon>
            </v-btn>
          </template>

          <v-list height="100%">
            <v-list-item class="d-flex flex-row justify-center pa-0 mt-n1 px-3">
              <v-subheader class="text-body-2 pa-0">Unit Type</v-subheader>
              <v-spacer/>
              <v-switch
                  v-model="showConversionToggle"
                  class="ml-4 my-auto"
                  color="primary"
                  dense
                  hide-details
                  inset>
                <template #prepend>
                  <span class="text-caption grey--text text--lighten-1">kW</span>
                </template>
                <template #append>
                  <span class="text-caption grey--text text--lighten-1 ml-n4">CO₂</span>
                </template>
              </v-switch>
            </v-list-item>
            <v-list-item class="pa-0 d-flex flex-row justify-center px-3" dense>
              <v-subheader class="text-body-2 pa-0">Duration</v-subheader>
              <v-spacer/>
              <v-btn-toggle
                  v-model="activeDuration"
                  active-class="primary--text"
                  color="transparent"
                  dense
                  mandatory>
                <v-btn
                    v-for="option in durationOptions"
                    class="transparent grey--text text--lighten-1"
                    :key="option.text"
                    small
                    :value="option.value">
                  <span class="hidden-sm-and-down">{{ option.text }}</span>
                </v-btn>
              </v-btn-toggle>
            </v-list-item>
          </v-list>
        </v-menu>
      </div>
    </template>
    <span>Options</span>
  </v-tooltip>
</template>

<script setup>
import {computed} from 'vue';
import {DAY, HOUR, MINUTE} from '@/components/now';

const props = defineProps({
  durationOption: {
    type: Object,
    default: () => {
    }
  },
  showConversion: {
    type: Boolean,
    default: false
  }
});
const emits = defineEmits(['update:durationOption', 'update:showConversion']);

// Defining the options for the duration type buttons
const durationOptions = [
  {
    text: '24H',
    value: {
      span: 20 * MINUTE,
      timeFrame: 24 * HOUR
    }
  },
  {
    text: '1W',
    value: {
      span: 1 * HOUR,
      timeFrame: 7 * DAY
    }
  },
  {
    text: '30D',
    value: {
      span: 1 * DAY,
      timeFrame: 30 * DAY
    }
  }
];

// Computed property to toggle between kW and CO2
// Syncs with the parent component
const showConversionToggle = computed({
  get: () => props.showConversion,
  set: (value) => emits('update:showConversion', value)
});

// Computed property to toggle between the duration options
// Syncs with the parent component
const activeDuration = computed({
  get: () => props.durationOption,
  set: (value) => emits('update:durationOption', value)
});
</script>
