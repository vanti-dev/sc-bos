<template>
  <v-card>
    <v-card-title>
      <v-expand-x-transition>
        <v-btn rounded="circle" v-if="pageIsSelected" @click="activeFilter = null" class="ml-n2 mr-2">
          <v-icon>mdi-arrow-left</v-icon>
        </v-btn>
      </v-expand-x-transition>
      <v-text-field v-bind="topTextInputBind" clearable hide-details density="compact" v-model="topSearch"/>
    </v-card-title>
    <v-card-subtitle class="pt-2 pb-0">
      <v-btn variant="text" block size="small" @click="topClear()" :disabled="!topIsDefaultChoice" class="overlap">
        <v-slide-x-transition>
          <span v-if="!pageIsSelected">Clear All Filters</span>
        </v-slide-x-transition>
        <v-slide-x-reverse-transition>
          <span v-if="pageIsSelected">Clear {{ pageFilter.title }}</span>
        </v-slide-x-reverse-transition>
      </v-btn>
    </v-card-subtitle>
    <v-window v-model="pageIndex">
      <v-window-item :value="1">
        <v-list>
          <v-fade-transition group>
            <template v-for="filter in displayFilters">
              <boolean-choice-list-item
                  v-if="filter.type === 'boolean'"
                  :key="filter.key"
                  :title="filter.title"
                  :icon="filter.icon"
                  :choice="choices[filter.key]"
                  :default-choice="isDefaultChoice(choices[filter.key])"
                  @input="choose(filter.key, $event)"
                  @clear="clear(filter.key)"
                  v-bind="filterMenuSelect"/>
              <page-choice-list-item
                  v-else
                  :key="filter.key"
                  :filter="filter"
                  :choice="choices[filter.key]"
                  :default-choice="isDefaultChoice(choices[filter.key])"
                  @click="activeFilter = filter.key"
                  @clear="clear(filter.key)"
                  v-bind="filterMenuSelect"/>
            </template>
          </v-fade-transition>
        </v-list>
      </v-window-item>
      <v-window-item :value="2">
        <list-chooser
            v-if="pageType === 'list'"
            :items="pageItems"
            :value="pageChoice?.value"
            @input="pageChooseOrBack($event)"/>
        <range-chooser
            v-else-if="pageType === 'range'"
            :items="pageItems"
            :value="pageChoice?.value"
            @input="pageChoose($event)"/>
        <span v-else>No filter selected</span>
      </v-window-item>
    </v-window>
  </v-card>
</template>
<script setup>
import BooleanChoiceListItem from '@/components/filter/BooleanChoiceListItem.vue';
import {filterCtxSymbol} from '@/components/filter/filterCtx.js';
import ListChooser from '@/components/filter/ListChooser.vue';
import PageChoiceListItem from '@/components/filter/PageChoiceListItem.vue';
import usePageCtx from '@/components/filter/pageCtx.js';
import RangeChooser from '@/components/filter/RangeChooser.vue';
import {computed, inject, provide, ref, watch} from 'vue';

const props = defineProps({
  ctx: {
    type: Object, // import('./filterCtx.js')
    default: () => ({})
  }
});

const ctx = /** @type {FilterCtx} */ inject(filterCtxSymbol, () => props.ctx, true);
provide(filterCtxSymbol, ctx);
const {
  active,
  filterSearch, displayFilters,
  choices, choose,
  clear,
  hasNonDefaultChoices,
  isDefaultChoice
} = ctx;

const activeFilter = ref(/** @type {FilterKey | null} */ null);
const {
  filter: pageFilter,
  type: pageType,
  title: pageTitle,
  choice: pageChoice,
  choose: pageChoose,
  clear: pageClear,
  isSelected: pageIsSelected,
  isNonDefaultChoice: pageIsNonDefaultChoice,
  items: pageItems,
  search: pageSearch
} = usePageCtx(ctx, activeFilter);

watch(active, (value, oldValue) => {
  if (value && !oldValue) {
    activeFilter.value = null;
    filterSearch.value = '';
    pageSearch.value = '';
  }
});

const topClear = () => {
  if (pageIsSelected.value) pageClear();
  else clear();
};
const topIsDefaultChoice = computed(() => {
  if (pageIsSelected.value) return pageIsNonDefaultChoice.value;
  return hasNonDefaultChoices.value;
});
const topTextInputBind = computed(() => {
  const allowInput = pageType.value !== 'range';
  return {
    placeholder: topPlaceholder.value,
    outlined: true,
    readonly: !allowInput,
    disabled: !allowInput,
    autofocus: allowInput
  };
});
const topPlaceholder = computed(() => {
  if (!pageIsSelected.value) return 'Filter by...';
  if (pageType.value === 'range') return `Adjust ${pageTitle.value.toLowerCase()} range`;
  return `Choose a ${pageTitle.value.toLowerCase()}`;
});

// Things we do to make keyboard navigation of the items work when in a v-menu.
// Must be added to all v-list-item on the filter select window.
const filterMenuSelect = computed(() => {
  const bind = {};
  // This "hack" works around VMenu behaviour of considering all v-list-item children that are part of the menu content
  // as part of it's up/down key navigation.
  // A tabindex of -1 triggers the exclusion behaviour so when on filter pages the items on the first page aren't
  // navigable with the up/down keys.
  // The reason we hit this issue is we're using v-window for transitions and that keeps everything in the DOM.
  if (activeFilter.value !== null) {
    bind.tabindex = -1;
  } else {
    // note: we can't just omit it, something remembers it and breaks the exclusion behaviour
    bind.tabindex = 0;
  }
  return bind;
});

// UX feature that allows selecting the same value again to go back to the filter choice page.
const pageChooseOrBack = (value) => {
  const changed = pageChoose(value);
  if (!changed) activeFilter.value = null;
};

// search text binding for the common text input between filter selection and list choice pages.
const topSearch = computed({
  get: () => {
    if (pageIsSelected.value) return pageSearch.value;
    return filterSearch.value;
  },
  set(value) {
    if (pageIsSelected.value) pageSearch.value = value;
    else filterSearch.value = value;
  }
});
const pageIndex = computed({
  get: () => pageIsSelected.value ? 2 : 1,
  set: (value) => {
    if (value === 1) activeFilter.value = null;
  }
});
</script>

<style scoped>
.v-input--switch.indeterminate ::v-deep(.v-input--switch__thumb),
.v-input--switch.indeterminate ::v-deep(.v-input--selection-controls__ripple) {
  transform: translate(10px, 0) scale(0.5) !important;
}

.v-btn.overlap ::v-deep(.v-btn__content) {
  display: grid;
}

.v-btn.overlap ::v-deep(.v-btn__content) > * {
  grid-row: 1;
  grid-column: 1;
}
</style>
