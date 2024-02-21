<template>
  <StatusAlert v-if="props.streamError" :resource="props.streamError"/>

  <span v-else-if="value && !props.streamError" class="text-no-wrap el-cell">
    <v-tooltip bottom>
      <template #activator="{ on, attrs }">
        <span v-on="on" v-bind="attrs">
          <v-icon :left="hasTotals" :class="{justEntered}" size="20">mdi-location-enter</v-icon>
          <span :class="{justEntered}" v-if="hasTotals">{{ enterTotal }}</span>
        </span>
      </template>
      Entered
    </v-tooltip>
    <v-divider vertical class="mx-2"/>
    <v-tooltip bottom>
      <template #activator="{ on, attrs }">
        <span v-on="on" v-bind="attrs">
          <span :class="{justLeft}" v-if="hasTotals">{{ leaveTotal }}</span>
          <v-icon :right="hasTotals" :class="{justLeft}" size="20">mdi-location-exit</v-icon>
        </span>
      </template>
      Left
    </v-tooltip>
  </span>
</template>
<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import {EnterLeaveEvent} from '@smart-core-os/sc-api-grpc-web/traits/enter_leave_sensor_pb';
import {computed, ref, watch} from 'vue';

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

const hasTotals = computed(() => props.value?.enterTotal !== undefined || props.value?.leaveTotal !== undefined);
const enterTotal = computed(() => props.value?.enterTotal || 0);
const leaveTotal = computed(() => props.value?.leaveTotal || 0);

const enterTimeoutHandle = ref(0);
const leaveTimeoutHandle = ref(0);
watch(() => props.value, (newVal, oldVal) => {
  if (!oldVal || !newVal) {
    justEntered.value = false;
    justLeft.value = false;
    clearTimeout(enterTimeoutHandle.value);
    clearTimeout(leaveTimeoutHandle.value);
    return;
  }

  if (newVal.direction === EnterLeaveEvent.Direction.ENTER) {
    justEntered.value = true;
    clearTimeout(enterTimeoutHandle.value);
    enterTimeoutHandle.value = setTimeout(() => {
      justEntered.value = false;
    }, props.showChangeDuration);
  }
  if (newVal.direction === EnterLeaveEvent.Direction.LEAVE) {
    justLeft.value = true;
    clearTimeout(leaveTimeoutHandle.value);
    leaveTimeoutHandle.value = setTimeout(() => {
      justLeft.value = false;
    }, props.showChangeDuration);
  }
}, {deep: true});
const justEntered = ref(false);
const justLeft = ref(false);
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
  color: var(--v-success-base);
}
</style>
