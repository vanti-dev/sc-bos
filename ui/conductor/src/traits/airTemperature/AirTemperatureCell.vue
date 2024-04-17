<template>
  <status-alert v-if="props.streamError" icon="mdi-thermometer-low" :resource="props.streamError"/>

  <span class="text-no-wrap at-cell" v-else-if="(hasTemp || hasSetPoint) && !props.streamError">
    <v-tooltip bottom v-if="hasTemp" open-delay="1000">
      <template #activator="{on, attrs}">
        <span v-bind="attrs" v-on="on">{{ tempStr }}</span>
      </template>
      <span>Current temperature</span>
    </v-tooltip>
    <v-icon class="mx-n1" v-if="hasTemp && hasSetPoint" size="20">mdi-menu-right</v-icon>
    <v-tooltip bottom v-if="hasSetPoint" open-delay="1000">
      <template #activator="{on, attrs}">
        <span v-bind="attrs" v-on="on">{{ setPointStr }}</span>
      </template>
      <span>Set point</span>
    </v-tooltip>
    <v-icon right size="20">mdi-thermometer</v-icon>
  </span>
</template>
<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import {useAirTemperature} from '@/traits/airTemperature/airTemperature.js';

const props = defineProps({
  value: {
    type: Object,
    default: () => {
    }
  },
  loading: {
    type: Boolean,
    default: false
  },
  showChangeDuration: {
    type: Number,
    default: 30 * 1000
  },
  streamError: {
    type: Object,
    default: null
  }
});

const {
  hasSetPoint,
  setPointStr,
  hasTemp,
  tempStr
} = useAirTemperature(() => props.value);
</script>

<style scoped>
.at-cell {
  display: flex;
  align-items: center;
}

.popup > .v-card__text {
  display: flex;
  justify-content: space-between;
  gap: 8px;
}
</style>
