<template>
  <status-alert v-if="props.streamError" icon="mdi-meter-electric-outline" :resource="props.streamError"/>

  <span class="text-no-wrap ed-cell" v-else-if="powerUseStr && !props.streamError">
    <v-tooltip location="bottom">
      <template #activator="{ props: _props }">
        <span v-bind="_props">
          <span>{{ powerUseStr }}</span>
          <v-icon end size="20">mdi-meter-electric-outline</v-icon>
        </span>
      </template>
      Power use
    </v-tooltip>
  </span>
</template>
<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import {useElectricDemand} from '@/traits/electricDemand/electric.js';
import {format} from '@/util/number.js';
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object, // of ElectricDemand.AsObject
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

const {realPower, realPowerUnit} = useElectricDemand(() => props.value);
const powerUseStr = computed(() => format(realPower.value, realPowerUnit.value));
</script>

<style scoped>
.el-cell {
  display: flex;
  align-items: center;
}
</style>
