<template>
  <v-col cols="auto" class="d-flex flex-row align-center py-0">
    <FilteredChips
        :active-filters="activeFilters"
        :active-type="activeType"
        :available-sources="availableSources"/>
    <FilterListMenu
        :active-type="activeType"
        :available-sources="availableSources"/>
  </v-col>
</template>
<script setup>
import FilteredChips from '@/components/filterBy/FilteredChips.vue';
import FilterListMenu from '@/components/filterBy/FilterListMenu.vue';
import useFilterBy from '@/components/filterBy/useFilterBy.js';
import {computed} from 'vue';

const props = defineProps({
  notification: {
    type: Boolean,
    default: false
  }
});

// Compute the active filter type from the props, so we can use it to extract the active filter store
const activeType = computed(() => {
  // eslint-disable-next-line no-unused-vars
  return Object.entries(props).find(([_, value]) => value === true)[0];
});

// Extract the active filter store
const {activeFilters, availableSources} = useFilterBy(() => activeType.value);
</script>
