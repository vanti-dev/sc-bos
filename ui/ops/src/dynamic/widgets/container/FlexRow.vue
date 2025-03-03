<template>
  <v-card :color="props.color" :title="props.title" :variant="variant" :style="{padding}">
    <div class="d-flex flex-row">
      <component
          v-for="(widget, index) in props.items"
          :key="widget.id ?? index"
          :is="widget.component"
          v-bind="widget.props"/>
    </div>
  </v-card>
</template>

<script setup>
import {computed} from 'vue';

const props = defineProps({
  items: {
    type: Array, // of widgets
    default: null
  },
  gap: {
    type: String,
    default: '24px'
  },
  justify: {
    type: String,
    default: 'flex-start'
  },
  align: {
    type: String,
    default: 'stretch'
  },
  wrap: {
    type: String,
    default: 'wrap'
  },
  itemMinWidth: {
    type: String,
    default: 'initial'
  },
  color: {
    type: String,
    default: ''
  },
  variant: {
    type: String,
    default: 'text'
  },
  title: {
    type: String,
    default: undefined
  }
});

const padding = computed(() => {
  // props that cause visible outline for the containing card
  if (!['', 'transparent'].includes(props.color) || ['elevated', 'outlined', 'tonal'].includes(props.variant)) {
    return '0 1em 1em';
  }
  return '0';
})
</script>

<style scoped>
.flex-row {
  gap: v-bind(gap);
  justify-content: v-bind(justify);
  align-items: v-bind(align);
  flex-wrap: v-bind(wrap);
}

.flex-row > * {
  min-width: v-bind(itemMinWidth);
}
</style>