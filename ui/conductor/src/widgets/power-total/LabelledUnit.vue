<template>
  <span class="labelled-unit">
    <span class="value">{{ valueStr }}</span>
    <span class="unit">{{ unitStr }}</span>
    <span class="label" :class="labelColor">{{ label }}</span>
  </span>
</template>

<script setup>
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Number,
    default: 0
  },
  unit: {
    type: String,
    default: 'kW'
  },
  label: {
    type: String,
    default: ''
  },
  labelColor: {
    type: String,
    default: null
  }
});

const valueStr = computed(() => {
  if (props.value === 0) return '0';
  return props.value?.toFixed(2) ?? '-';
});
const unitStr = computed(() => props.unit);
</script>

<style scoped>
.labelled-unit {
  display: inline-grid;
  grid-template-areas:
    "value unit"
    "label label";
  grid-template-columns: auto 1fr;

  line-height: 1;
  gap: 0 .06em;
}

.value {
  grid-area: value;
  align-self: baseline;
}
.unit {
  grid-area: unit;
  align-self: baseline;
  font-size: 50%;
  font-weight: lighter;
}
.label {
  grid-area: label;
  font-size: 40%;
  font-weight: lighter;
  letter-spacing: 1px;
}

</style>
