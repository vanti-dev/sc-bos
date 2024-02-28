<template>
  <v-list v-if="selectedOption === null" class="px-0 mx-0 mt-5 overflow-y-auto" dense max-height="300">
    <v-list-item
        v-for="source in props.availableSources"
        active-class="rounded"
        class="ml-0 text-truncate rounded"
        :key="source.title"
        :title="source.title"
        @click="selectedOption = source">
      <v-list-item-icon class="ml-n2">
        <v-icon size="22">
          {{ source.icon }}
        </v-icon>
      </v-list-item-icon>
      <v-list-item-content>
        <v-list-item-title class="text-body-2 d-flex flex-row justify-space-between align-center">
          {{ source.title }}
          <v-icon v-if="findActiveFilter(source.title)" class="ml-2" color="primary" size="16">
            mdi-tag-outline
          </v-icon>
          <v-icon class="ml-auto mr-0">
            mdi-chevron-right
          </v-icon>
        </v-list-item-title>
      </v-list-item-content>
    </v-list-item>
  </v-list>

  <v-list v-else class="px-0 mx-0 mt-4 rounded-0 pb-1 overflow-y-auto" dense max-height="300">
    <div v-if="selectedOption?.type === 'range'">
      <v-subheader class="mt-n4 mb-2 ml-n2 text-caption sticky-wrapper neutral">
        Severity not above
      </v-subheader>
      <v-select
          class="ml-0"
          dense
          flat
          hide-details
          :items="selectedOption.value"
          item-text="label"
          item-value="value"
          label="Select a range"
          outlined
          solo
          style="width: 100%"
          :value="findActiveFilter('severityNotAbove')"
          @change="updateFilter('severityNotAbove', $event)"/>
      <v-subheader class="mt-3 mb-1 ml-n2 text-caption sticky-wrapper neutral">
        Severity not below
      </v-subheader>
      <v-select
          class="ml-0 mb-2"
          dense
          flat
          hide-details
          :items="selectedOption.value"
          item-text="label"
          item-value="value"
          label="Select a range"
          outlined
          solo
          style="width: 100%"
          :value="findActiveFilter('severityNotBelow')"
          @change="updateFilter('severityNotBelow', $event)"/>
    </div>
    <div v-else>
      <v-subheader class="mt-n4 mb-1 ml-n2 text-caption sticky-wrapper neutral">
        {{ selectedOption.title }}
      </v-subheader>
      <v-list-item
          v-for="value in props.searchedSelectedOption"
          :key="value"
          active-class="rounded"
          class="ml-0 text-truncate rounded"
          :title="value"
          @click="updateFilter(selectedOption.title, value)">
        <v-list-item-content class="ml-n1">
          <v-list-item-title class="text-body-2 d-flex flex-row align-center">
            {{ value }}
            <v-icon v-if="findActiveFilter(value)" class="ml-auto mr-0 float-right" color="primary" size="18">
              mdi-tag-outline
            </v-icon>
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </div>
  </v-list>
</template>
<script setup>
import {computed} from 'vue';
import useFilterBy from '@/components/filterBy/useFilterBy.js';

const emits = defineEmits(['update:selectedOption']);
const props = defineProps({
  activeType: {
    type: String,
    default: ''
  },
  availableSources: {
    type: Array,
    default: () => []
  },
  filterInputValue: {
    type: String,
    default: ''
  },
  searchedSelectedOption: {
    type: Array,
    default: () => []
  },
  selectedOption: {
    type: Object,
    default: null
  }
});

const {activeFilters, updateFilter} = useFilterBy(() => props.activeType);

const findActiveFilter = (value) => {
  const isSeverity = activeFilters.value.some((filter) => {
    const filterKeyLower = filter.key.toLowerCase();
    return filterKeyLower === 'severitynotbelow' || filterKeyLower === 'severitynotabove';
  });

  if (value === 'Severity' && isSeverity) {
    return {title: 'Severity'};
  }

  return activeFilters.value.find((filter) => {
    // Convert filter.value to string to safely call toLowerCase
    const filterValueStr = String(filter.value).toLowerCase();
    const filterKeyLower = filter.key.toLowerCase();
    return filterValueStr === value.toLowerCase() || filterKeyLower === value.toLowerCase();
  });
};

const selectedOption = computed({
  get: () => props.selectedOption,
  set: (value) => emits('update:selectedOption', value)
});
</script>

<style scoped>
.sticky-wrapper {
  position: sticky;
  top: -11px; /* Adjust this value based on your needs */
  z-index: 1000; /* Ensure the sticky element stays above other content */
}
</style>

