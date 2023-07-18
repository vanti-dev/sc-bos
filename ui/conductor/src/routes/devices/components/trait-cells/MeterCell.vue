<template>
  <span class="text-no-wrap ed-cell" v-if="meterReading">
    <span>{{ meterReading }}</span>
    <v-icon right size="20">mdi-counter</v-icon>
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
  type: {
    type: Object,
    default: () => {
    }
  },
  loading: {
    type: Boolean,
    default: false
  }

});

const meterReading = computed(() => {
  return `${props.value?.value?.usage.toFixed(2)} ${props.type?.response?.unit}`;
});
</script>

<style scoped>
.el-cell {
  display: flex;
  align-items: center;
}
</style>
