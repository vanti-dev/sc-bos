<template>
  <span class="text-no-wrap ed-cell" v-if="meterReading">
    <span>{{ meterReading }} kWh</span>
  </span>
</template>
<script setup>
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object, // of MeterReading.AsObject
    default: () => {
    }
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const meterReading = computed(() => {
  return props.value?.usage.toFixed(2);
});
</script>

<style scoped>
.el-cell {
  display: flex;
  align-items: center;
}
</style>
