<template>
  <content-card class="mt-3 mb-3">
    <div class="d-flex flex-row mb-2">
      <v-card-title class="text-h4 pl-4">Presence</v-card-title>
      <StatusAlert :resource="occupancyValue.streamError"/>
    </div>
    <v-col cols="12" class="d-flex flex-column pl-4">
      <v-row class="d-flex flex-row align-center px-3 pb-2">
        <span v-if="occupancyValue.value" :class="[stateColor, 'text-h6 font-weight-bold']">
          {{ stateStr }}
        </span>
        <v-spacer/>
        <div class="d-flex flex-row ma-0 text-caption font-weight-regular">
          <span class="mr-1">Last presence detected:</span>
          <span>{{ timeAgo }}</span>
        </div>
      </v-row>
    </v-col>
  </content-card>
</template>

<script setup>
import {occupancyStateToString} from '@/api/sc/traits/occupancy';
import ContentCard from '@/components/ContentCard.vue';
import {DAY, HOUR, MINUTE, SECOND, useNow} from '@/components/now';
import StatusAlert from '@/components/StatusAlert.vue';
import useOccupancyTrait from '@/composables/traits/useOccupancyTrait';
import {formatTimeAgo} from '@/util/date';
import {Occupancy} from '@smart-core-os/sc-api-grpc-web/traits/occupancy_sensor_pb';
import {computed, ref, watch} from 'vue';

const props = defineProps({
  name: {
    type: String,
    required: true
  }
});

const {occupancyValue} = useOccupancyTrait(props);

const state = computed(() => {
  return occupancyValue?.value?.state;
});
const stateStr = computed(() => {
  if (state.value === undefined) return '';
  return occupancyStateToString(state.value);
});

const stateColor = computed(() => {
  if (state.value === Occupancy.State.OCCUPIED) {
    return 'success--text text--lighten-1';
  } else if (state.value === Occupancy.State.UNOCCUPIED) {
    return 'warning--text';
  } else if (state.value === Occupancy.State.IDLE) {
    return 'info--text';
  } else {
    return undefined;
  }
});

// Create a lastChecked timestamp (for second to be used in the status popup
const {now} = useNow(SECOND);
const lastChecked = ref(null);

// Update lastChecked timestamp when isLoading changes
watch(() => occupancyValue?.updateTime, (updated) => {
  lastChecked.value = Date.parse(updated);
}, {immediate: true});

// Create a timeAgo computed property to display time in words
const timeAgo = computed(() => {
  if (!lastChecked.value) return 'Never';
  return formatTimeAgo(lastChecked.value, now.value, MINUTE, HOUR, DAY);
});
</script>
