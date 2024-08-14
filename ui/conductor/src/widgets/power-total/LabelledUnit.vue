<template>
  <span class="labelled-unit">
    <template v-if="showErr">
      <v-tooltip bottom>
        <template #activator="{ props }">
          <span style="height: 1em" class="err">
            <v-icon v-bind="props" color="error" size=".75em">mdi-alert-circle-outline</v-icon>
          </span>
        </template>
        <span>{{ errStr }}</span>
      </v-tooltip>
    </template>
    <span class="value" v-if="!showErr">{{ valueStr }}</span>
    <span class="unit" v-if="!showErr">{{ unitStr }}</span>
    <span class="label" :class="labelColor">{{ label }}</span>
  </span>
</template>

<script setup>
import useError from '@/composables/error.js';
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
  },
  error: {
    type: [Object, String],
    default: null
  }
});

const valueStr = computed(() => {
  if (props.value === 0) return '0';
  return props.value?.toFixed(2) ?? '-';
});
const unitStr = computed(() => props.unit);

const {errStr, showErr} = useError(() => props.error);
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

.err {
  grid-column: 1 / -1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

</style>
