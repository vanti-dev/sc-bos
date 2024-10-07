<template>
  <status-alert v-if="props.streamError" icon="mdi-thermometer-low" :resource="props.streamError"/>

  <air-temperature-chip
      v-else-if="(hasTemp || hasSetPoint) && !props.streamError"
      variant="text" size="30" layout="right"
      :current-temp="temp" :set-point="setPoint"/>
</template>
<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import {useAirTemperature} from '@/traits/airTemperature/airTemperature.js';
import AirTemperatureChip from '@/traits/airTemperature/AirTemperatureChip.vue';

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
  hasTemp,
  setPoint,
  temp
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
