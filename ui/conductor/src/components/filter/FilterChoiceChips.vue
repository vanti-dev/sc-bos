<template>
  <span class="chip-list">
    <template v-for="choice in chipChoices">
      <boolean-chooser-chip
          v-if="choiceType(choice) === 'boolean'"
          :key="choice.filter"
          :ctx="usePageCtx(ctx, choice.filter)"
          color="neutral-lighten-2"/>
      <menu-chooser-chip
          v-else-if="choiceType(choice) === 'list'"
          :key="choice.filter"
          :ctx="usePageCtx(ctx, choice.filter)"
          color="neutral-lighten-2"
          @active="activateChip($event, choice)"
          v-slot="{items, value, choose}">
        <list-chooser :items="items" :value="value" @input="choose"/>
      </menu-chooser-chip>
      <menu-chooser-chip
          v-else-if="choiceType(choice) === 'range'"
          :key="choice.filter"
          :ctx="usePageCtx(ctx, choice.filter)"
          color="neutral-lighten-2"
          @active="activateChip($event, choice)"
          v-slot="{items, value, choose}">
        <range-chooser :items="items" :value="value" @input="choose"/>
      </menu-chooser-chip>
      <v-chip
          v-else
          :key="choice.filter"
          @click:close="clear(choice.filter)"
          closable
          color="neutral-lighten-2">
        {{ choice.text ?? choice.value }}
      </v-chip>
    </template>
  </span>
</template>

<script setup>

import BooleanChooserChip from '@/components/filter/BooleanChooserChip.vue';
import {filterCtxSymbol} from '@/components/filter/filterCtx.js';
import ListChooser from '@/components/filter/ListChooser.vue';
import MenuChooserChip from '@/components/filter/MenuChooserChip.vue';
import usePageCtx from '@/components/filter/pageCtx.js';
import RangeChooser from '@/components/filter/RangeChooser.vue';
import {computed, inject, ref} from 'vue';

const props = defineProps({
  ctx: {
    type: Object, // of import('./filterCtx.js').FilterCtx,
    default: null
  }
});

const ctx = /** @type {FilterCtx} */ inject(filterCtxSymbol, () => props.ctx, true);
const {sortedChoices, nonDefaultChoices, clear, filtersByKey} = ctx;
const choiceType = (choice) => filtersByKey.value[choice.filter].type;

const activateChip = (b, choice) => {
  if (b) {
    activeChip.value = choice;
  } else if (activeChip.value?.filter === choice.filter) {
    activeChip.value = null;
  }
};
const activeChip = ref(null);

// equivalent to nonDefaultChoices + activeChip
const chipChoices = computed(() => {
  const ac = activeChip.value;
  const ndc = nonDefaultChoices.value;
  if (ac === null) {
    return ndc;
  }

  const res = [];
  let i = 0;
  for (const sc of sortedChoices.value) {
    if (ndc.length > i && ndc[i].filter === sc.filter) {
      i++;
      res.push(ndc);
      continue;
    }
    if (ac.filter === sc.filter) {
      res.push(ac);
    }
  }
  return res;
});

</script>

<style scoped>
.chip-list {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin: 0 -4px;
}
</style>
