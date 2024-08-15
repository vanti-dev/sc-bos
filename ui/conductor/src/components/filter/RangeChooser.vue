<template>
  <div>
    <div class="text-subtitle-2 justify-center text-h6">{{ summaryStr }}</div>
    <v-range-slider
        :min="min"
        :max="max"
        :step="1"
        show-ticks="always"
        tick-size="4"
        v-model="model"
        hide-details
        class="mx-4 mb-4"
        prepend-icon="mdi-infinity"
        append-icon="mdi-infinity"/>
  </div>
</template>

<script setup>
import choiceRangeStr from '@/components/filter/choiceRange.js';
import {computed, ref, watch} from 'vue';

const props = defineProps({
  items: {
    type: Array, // of import('./filterCtx.js').FilterItem
    default: () => []
  },
  value: {
    type: Object, // import('./filterCtx.js').ChoiceRange
    default: () => null
  }
});
const emits = defineEmits(['input']);

// We want ticks for all items + 2, one before and one after.
// To make indexing easier we range from -1 to items.length
const min = computed(() => -1);
const max = computed(/** @type {() => number} */() => {
  if (props.items?.length > 0) return props.items.length;
  return 0;
});

const itemValue = (item) => item?.value ?? item;
const valueToStep = (item, def = 0) => {
  const val = itemValue(item);
  if (val === null || val === undefined) return def;
  const i = props.items.findIndex((v) => itemValue(v) === val);
  if (i === -1) return def;
  return i;
};
const valueToModel = (val) => {
  return [
    valueToStep(val?.from, min.value),
    valueToStep(val?.to, max.value)
  ];
};
const modelToValue = ([from, to]) => {
  if (from === min.value && to === max.value) return undefined; // aka all
  const res = {};
  if (from !== min.value) res.from = props.items[from];
  if (to !== max.value) res.to = props.items[to];
  return res;
};

const modelVal = ref(/** @type {null | [number,number]} */ null);
watch(() => props.value, () => {
  modelVal.value = null;
});
const model = computed({
  get: () => {
    if (modelVal.value !== null) return modelVal.value;
    return valueToModel(props.value);
  },
  set: (val) => {
    // check that both ticks are not min or max, this is the same as all.
    if (val[0] === min.value && val[1] === min.value) val[1]++;
    if (val[0] === max.value && val[1] === max.value) val[0]--;
    modelVal.value = val;
    emits('input', modelToValue(val));
  }
});
const localValue = computed(() => modelToValue(model.value));

const summaryStr = computed(() => choiceRangeStr(localValue.value));
</script>

<style scoped>

</style>
