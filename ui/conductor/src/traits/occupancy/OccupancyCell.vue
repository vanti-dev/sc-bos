<template>
  <StatusAlert v-if="props.streamError" icon="mdi-crosshairs" :resource="props.streamError"/>

  <v-menu
      v-else
      left
      bottom
      offset-y
      nudge-bottom="4px"
      nudge-right="4px"
      transition="slide-x-reverse-transition"
      open-on-hover>
    <template #activator="{on}">
      <v-icon :class="state" :color="iconColor" v-on="on" size="20">{{ icon }}</v-icon>
    </template>
    <v-card :color="stateColor">
      <v-card-text class="py-2">{{ tooltipStr }}</v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import {useTimeSince} from '@/composables/time.js';
import {useOccupancy} from '@/traits/occupancy/occupancy.js';
import {computed} from 'vue';

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
  streamError: {
    type: Object,
    default: null
  }
});

const {state, stateStr, stateColor, icon, iconColor, lastUpdate} = useOccupancy(() => props.value);
const {showTimeSince, timeSinceStr} = useTimeSince(lastUpdate);

const tooltipStr = computed(() => {
  let s = stateStr.value;
  if (showTimeSince.value) {
    s += ` for ${timeSinceStr.value}`;
  }
  return s;
});
</script>
