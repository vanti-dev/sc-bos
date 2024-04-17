<template>
  <status-alert v-if="props.streamError" icon="mdi-counter" :resource="props.streamError"/>

  <span class="text-no-wrap ed-cell" v-else-if="value && !props.streamError">
    <v-tooltip bottom>
      <template #activator="{ on, attrs }">
        <span v-on="on" v-bind="attrs">
          <span>{{ usageAndUnit }}</span>
          <v-icon right size="20">mdi-counter</v-icon>
        </span>
      </template>
      <span>Meter reading</span>
    </v-tooltip>
  </span>
</template>
<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import {useMeterReading} from '@/traits/meter/meter.js';

const props = defineProps({
  value: {
    type: Object, // of MeterReading.AsObject
    default: () => {
    }
  },
  info: {
    type: Object, // of MeterReadingInfo.AsObject
    default: () => null
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

const {usageAndUnit} = useMeterReading(() => props.value, () => props.info);
</script>

<style scoped>
.el-cell {
  display: flex;
  align-items: center;
}
</style>
