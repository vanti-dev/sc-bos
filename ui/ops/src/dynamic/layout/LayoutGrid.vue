<template>
  <div class="grid--container" :class="{signage: signageEnabled}" :style="signageStyles">
    <component
        :is="cellComponent(cell)"
        v-for="(cell, index) in props.cells"
        :key="cell.id ?? cell.component ?? index"
        :style="cellStyles(cell)"
        v-bind="cell.props"
        class="grid--cell"/>
  </div>
</template>

<script setup>
import useSignage from '@/composables/signage.js';
import PlaceholderCard from '@/dynamic/widgets/general/PlaceholderCard.vue';
import {computed} from 'vue';

const props = defineProps({
  cells: {
    type: Array,
    required: true
  }
});
const countLines = (cells, inlineStartProp, inlineSpanProp) => {
  let lines = 0;
  for (const cell of cells) {
    const end = cell.loc[inlineStartProp] + cell.loc[inlineSpanProp];
    if (end > lines) {
      lines = end;
    }
  }
  return lines;
}
// these have -1 because there's 1 less column/row than lines: | col1 | col2 | <- cols = 2, lines = 3
const columnCount = computed(() => countLines(props.cells, 'x', 'w') - 1);
const rowCount = computed(() => countLines(props.cells, 'y', 'h') - 1);

const {enabled: signageEnabled, styles: signageStyles} = useSignage();

const cellStyles = (cell) => {
  return {
    '--x': cell.loc.x,
    '--y': cell.loc.y,
    '--w': cell.loc.w,
    '--h': cell.loc.h,
  };
}
const cellComponent = (cell) => {
  return cell.component ?? PlaceholderCard;
}
</script>

<style scoped lang="scss">
.grid--container {
  --gap: 10px;
  display: grid;
  grid-template-columns: repeat(v-bind(columnCount), 1fr);
  align-content: stretch;
  gap: var(--gap);
  min-height: 100%;
}

.signage.grid--container {
  grid-template-rows: repeat(v-bind(rowCount), 1fr);
  padding: var(--gap);
  // add scrollbar gutter to match the app one when needed
  scrollbar-gutter: stable;

  .grid--cell {
    overflow: hidden;
  }
}

.grid--cell {
  grid-column-start: var(--x);
  grid-column-end: span var(--w);
  grid-row-start: var(--y);
  grid-row-end: span var(--h);
  min-height: 0;
  min-width: 0;
  overflow: auto;
}
</style>