<template>
  <v-tooltip left nudge-right="8px">
    <template #activator="{attr, on}">
      <v-icon v-bind="attr" v-on="on" size="20" :color="colorClass">{{ iconStr }}</v-icon>
    </template>
    <span>{{ tooltipStr }}</span>
  </v-tooltip>
</template>
<script setup>
import {Emergency} from '@smart-core-os/sc-api-grpc-web/traits/emergency_pb';
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object, // of Emergency.AsObject
    default: () => {
    }
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const colorClass = computed(() => {
  const val = props.value?.level;
  const drill = props.value?.drill;
  switch (val) {
    default:
    case Emergency.Level.OK:
      return '';
    case Emergency.Level.WARNING:
      return 'warning--text';
    case Emergency.Level.EMERGENCY:
      if (drill) {
        return 'info--text';
      }
      return 'error--text';
  }
});
const iconStr = computed(() => {
  const val = props.value?.level;
  switch (val) {
    default:
    case Emergency.Level.OK:
      return 'mdi-smoke-detector-outline';
    case Emergency.Level.WARNING:
      return 'mdi-smoke-detector';
    case Emergency.Level.EMERGENCY:
      return 'mdi-smoke-detector-alert';
  }
});
const tooltipStr = computed(() => {
  // todo: work out a better message based on current state
  return 'Emergency status';
});
</script>

<style scoped>
.el-cell {
  display: flex;
  align-items: center;
}
</style>
