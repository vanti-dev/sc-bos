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
            <v-list-item
                class="d-flex flex-row justify-center pa-0 mt-n1 px-3"
                @click.stop="showConversionToggle = !showConversionToggle">
              <v-subheader class="text-body-2 pa-0">Unit Type</v-subheader>
              <v-spacer/>
              <v-switch
                  class="ml-4 my-auto no-pointer-events"
                  color="primary"
                  dense
                  readonly
                  hide-details
                  inset
                  :value="props.showConversion">
                <template #prepend>
                  <span class="text-caption grey--text text--lighten-1">kW</span>
                </template>
                <template #append>
                  <span class="text-caption grey--text text--lighten-1 ml-n4">CO₂</span>
                </template>
              </v-switch>
            </v-list-item>
            <v-list-item class="pa-0 d-flex flex-row justify-center px-3" dense @click.stop="changeDuration">
              <v-subheader class="text-body-2 pa-0">Duration</v-subheader>
              <v-spacer/>
              <v-btn-toggle
                  active-class="primary"
                  dense
                  :value="durationOption">
                <v-btn
                    v-for="option in durationOptions"
                    active-class="primary text--darken-3"
                    class="transparent grey--text text--lighten-1 no-pointer-events btn-no-hover"
                    :key="option.id"
                    :ripple="false"
                    small
                    :value="option.value">
                  <span class="text-caption">{{ option.text }}</span>
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
import {DAY, HOUR, MINUTE} from '@/components/now';
import {computed} from 'vue';

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

// Click handler on duration row.
// Increasing the duration on each click to the next level. If reaching the last option
// returns to the first available option and restarts the cycle.
const changeDuration = () => {
  const index = durationOptions.findIndex(option => option.value.id === activeDuration.value.id);
  const nextIndex = (index + 1) % durationOptions.length;
  activeDuration.value = durationOptions[nextIndex].value;
};

</script>

<style lang="scss" scoped>
.no-pointer-events {
  pointer-events: none;
}

.v-btn-toggle.no-pointer-events .v-btn {
  background-color: transparent !important;
  color: inherit !important; /* Adjust based on your needs */
}
</style>
