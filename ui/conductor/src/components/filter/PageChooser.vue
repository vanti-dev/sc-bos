<template>
  <v-card>
    <v-card-title>
      <v-text-field v-bind="topTextInputBind" clearable hide-details density="compact" v-model="search"/>
    </v-card-title>
    <slot/>
  </v-card>
</template>

<script setup>
import {computed} from 'vue';

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  type: {
    type: String,
    default: 'list'
  },
  search: {
    type: String,
    default: ''
  }
});
const emits = defineEmits(['update:search']);

const search = computed({
  get: () => props.search,
  set: (value) => emits('update:search', value)
});

const topTextInputBind = computed(() => {
  const allowInput = props.type !== 'range';
  return {
    placeholder: topPlaceholder.value,
    outlined: true,
    readonly: !allowInput,
    disabled: !allowInput,
    autofocus: allowInput
  };
});
const topPlaceholder = computed(() => {
  if (props.type === 'range') return `Adjust ${props.title.toLowerCase()} range`;
  return `Choose a ${props.title.toLowerCase()}`;
});
</script>
