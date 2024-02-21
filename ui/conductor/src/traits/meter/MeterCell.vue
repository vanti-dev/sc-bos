<template>
  <StatusAlert v-if="props.streamError" icon="mdi-counter" :resource="props.streamError"/>

  <span class="text-no-wrap ed-cell" v-else-if="meterReading && !props.streamError">
    <v-tooltip bottom>
      <template #activator="{ on, attrs }">
        <span v-on="on" v-bind="attrs">
          <span>{{ meterReading }}</span>
          <v-icon right size="20">mdi-counter</v-icon>
        </span>
      </template>
      <span>Meter reading</span>
    </v-tooltip>
  </span>
</template>
<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object, // of MeterReading.AsObject
    default: () => {
    }
  },
  unit: {
    type: String,
    default: ''
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

const meterReading = computed(() => {
  let val = props.value?.usage?.toFixed(2) ?? '';
  if (props.unit) {
    val += ` ${props.unit}`;
  }
  return val;
});
</script>

<style scoped>
.el-cell {
  display: flex;
  align-items: center;
}
</style>
