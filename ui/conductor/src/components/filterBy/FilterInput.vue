<template>
  <v-text-field
      clearable
      class="elevation-0"
      color="primary"
      dense
      flat
      :label="filterInputLabel"
      hide-details
      outlined
      single-line
      solo
      v-model="filterInputModel"/>
</template>

<script setup>
import {computed} from 'vue';

const emits = defineEmits(['update:filterInputValue']);
const props = defineProps({
  filterInputValue: {
    type: String,
    default: ''
  },
  label: {
    type: String,
    default: ''
  }
});

const filterInputLabel = computed(() => {
  if (!props.label) return 'Filter by...';
  else return 'Filter by ' + props.label.toLowerCase();
});

const filterInputModel = computed({
  get: () => props.filterInputValue,
  set: (value) => emits('update:filterInputValue', value)
});
</script>

<style scoped>

</style>
