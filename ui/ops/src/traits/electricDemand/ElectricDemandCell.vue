<template>
  <status-alert v-if="props.streamError" icon="mdi-meter-electric-outline" :resource="props.streamError"/>

  <span class="text-no-wrap ed-cell" v-else-if="powerUseStr && !props.streamError">
    <v-tooltip location="bottom">
      <template #activator="{ props: _props }">
        <span v-bind="_props">
          <span v-if="realPower">{{ powerUseStr }}</span>
          <span v-else-if="current">{{ currentStr }} </span>
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

const {realPower, realPowerUnit, current, currentUnit} = useElectricDemand(() => props.value);
const powerUseStr = computed(() => format(realPower.value, realPowerUnit.value));
const currentStr = computed(() => format(current.value, currentUnit.value));
</script>

<style scoped>
.el-cell {
  display: flex;
  align-items: center;
}
</style>
