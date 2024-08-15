<template>
  <status-alert v-if="props.streamError" icon="mdi-smoke-detector-outline" :resource="props.streamError"/>

  <v-tooltip v-else location="left" nudge-right="8px">
    <template #activator="{props}">
      <v-icon v-bind="props" size="20" :color="colorClass">{{ iconStr }}</v-icon>
    </template>
    <span>{{ tooltipStr }}</span>
  </v-tooltip>
</template>
<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import {useEmergency} from '@/traits/emergency/emergency.js';

const props = defineProps({
  value: {
    type: Object, // of Emergency.AsObject
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

const {colorClass, iconStr, tooltipStr} = useEmergency(() => props.value);
</script>

<style scoped>
.el-cell {
  display: flex;
  align-items: center;
}
</style>
