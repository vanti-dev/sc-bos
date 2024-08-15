<template>
  <v-tooltip v-if="needsTooltip" v-bind="$attrs" location="bottom">
    <template #activator="{props}">
      <v-chip v-bind="{...props, ...$attrs}">
        <slot>{{ _title }}</slot>
      </v-chip>
    </template>
    {{ props.name }}
  </v-tooltip>
  <v-chip v-else v-bind="$attrs">
    <slot>{{ _title }}</slot>
  </v-chip>
</template>

<script setup>
import {computed} from 'vue';

defineOptions({
  inheritAttrs: false
});

const props = defineProps({
  name: {
    type: String,
    default: ''
  },
  title: {
    type: String,
    default: ''
  }
});

const needsTooltip = computed(() => {
  return props.name !== _title.value;
});
const _title = computed(() => {
  if (props.title) {
    return props.title;
  }
  return props.name.split('/').pop();
});
</script>
