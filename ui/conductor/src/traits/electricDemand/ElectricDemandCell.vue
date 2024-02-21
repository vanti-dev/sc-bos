<template>
  <StatusAlert v-if="props.streamError" icon="mdi-meter-electric-outline" :resource="props.streamError"/>

  <span class="text-no-wrap ed-cell" v-else-if="powerUseStr && !props.streamError">
    <v-tooltip bottom>
      <template #activator="{ on, attrs }">
        <span v-on="on" v-bind="attrs">
          <span>{{ powerUseStr }}</span>
          <v-icon right size="20">mdi-meter-electric-outline</v-icon>
        </span>
      </template>
      Power use
    </v-tooltip>
  </span>
</template>
<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
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

const realPower = computed(() => {
  return props.value?.realPower;
});
const powerUseStr = computed(() => {
  if (!realPower.value) return '';
  return `${(realPower.value / 1000).toFixed(2)}kW`;
});
</script>

<style scoped>
.el-cell {
  display: flex;
  align-items: center;
}
</style>
