<template>
  <status-alert v-if="props.streamError" :resource="props.streamError"/>

  <span v-else-if="value && !props.streamError" class="text-no-wrap el-cell">
    <v-tooltip location="bottom">
      <template #activator="{ props }">
        <span v-bind="props">
          <v-icon :start="hasTotals" :class="{justEntered}" size="20">mdi-location-enter</v-icon>
          <span :class="{justEntered}" v-if="hasTotals">{{ enterTotal }}</span>
        </span>
      </template>
      Entered
    </v-tooltip>
    <v-divider vertical class="mx-2"/>
    <v-tooltip location="bottom">
      <template #activator="{ props }">
        <span v-bind="props">
          <span :class="{justLeft}" v-if="hasTotals">{{ leaveTotal }}</span>
          <v-icon :end="hasTotals" :class="{justLeft}" size="20">mdi-location-exit</v-icon>
        </span>
      </template>
      Left
    </v-tooltip>
  </span>
</template>
<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import {useEnterLeaveEvent} from '@/traits/enterLeave/enterLeave.js';

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
  hasTotals, enterTotal, leaveTotal,
  justEntered, justLeft
} = useEnterLeaveEvent(() => props.value, {
  showChangeDuration: () => props.showChangeDuration
});
</script>

<style scoped>
.el-cell {
  display: flex;
  align-items: center;
}

.el-cell > * {
  transition: color 0.2s ease-in-out;
}

.justEntered,
.justLeft {
  color: rgb(var(--v-theme-success));
}
</style>
