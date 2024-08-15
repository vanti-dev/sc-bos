<template>
  <v-switch
      v-bind="$attrs"
      :class="{indeterminate: _indeterminate}"
      :model-value="_value"
      @update:model-value="emit('input', $event)"
      v-on="$listeners">
    <template #label v-if="$slots.label">
      <slot name="label"/>
    </template>
  </v-switch>
</template>
<script setup>
import {computed, ref, watch} from 'vue';

const props = defineProps({
  value: {
    type: null,
    default: undefined
  },
  indeterminate: {
    type: Boolean,
    default: false
  }
});
const emit = defineEmits(['input', 'update:indeterminate']);

const _indeterminate = ref(/** @type {boolean} */ props.indeterminate ?? false);
watch(() => props.indeterminate, () => {
  _indeterminate.value = props.indeterminate;
});
watch(_indeterminate, (newValue, oldValue) => {
  if (newValue === oldValue) return;
  emit('update:indeterminate', newValue);
});

const _value = computed(() => _indeterminate.value ? false : props.value);

</script>
<style scoped>
.v-input--switch.indeterminate :deep(.v-input--switch__thumb),
.v-input--switch.indeterminate :deep(.v-input--selection-controls__ripple) {
  transform: translate(10px, 0) scale(0.5) !important;
}
</style>
