<template>
  <span class="text-no-wrap ed-cell" v-if="powerUseStr">
    <span>{{ powerUseStr }}</span>
    <v-icon right size="20">mdi-meter-electric-outline</v-icon>
  </span>
</template>
<script setup>
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
