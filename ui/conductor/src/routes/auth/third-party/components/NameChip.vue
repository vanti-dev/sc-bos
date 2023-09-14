<template>
  <v-tooltip v-if="needsTooltip" v-bind="$attrs" bottom>
    <template #activator="{on, attrs}">
      <v-chip v-on="on" v-bind="{...attrs, ...$attrs}">
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
