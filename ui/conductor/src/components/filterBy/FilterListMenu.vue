<template>
  <v-menu
      content-class="mt-2"
      :close-on-content-click="false"
      left
      max-width="290"
      min-width="290"
      nudge-bottom="0"
      nudge-left="0"
      nudge-right="0"
      nudge-top="0"
      offset-y
      v-model="showHideFilterList">
    <template #activator="{ on }">
      <FilterButton v-on="on" :show-filter-list.sync="showHideFilterList"/>
    </template>
    <ContentCard class="pb-0 pb-1 mb-0">
      <v-row class="px-2 pt-3 pb-0 align-center sticky-wrapper">
        <v-tooltip v-if="selectedOption" bottom>
          <template #activator="{ on }">
            <v-btn
                class="px-4 mr-1 rounded"
                height="40"
                icon
                v-on="on"
                width="22"
                @click="selectedOption = null">
              <v-icon>
                mdi-arrow-left
              </v-icon>
            </v-btn>
          </template>
          <span>Back</span>
        </v-tooltip>
        <FilterInput
            v-show="selectedOption?.type !== 'range'"
            :filter-input-value.sync="filterInputValue"
            :label="selectedOption?.title"/>
      </v-row>
      <FilterOptions
          :active-type="props.activeType"
          :available-sources="availableSearchedSources"
          :filter-input-value.sync="filterInputValue"
          :selected-option.sync="selectedOption"
          :searched-selected-option="searchedSelectedOption"/>
    </ContentCard>
  </v-menu>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import FilterInput from '@/components/filterBy/FilterInput.vue';
import FilterOptions from '@/components/filterBy/FilterOptions.vue';
import FilterButton from '@/components/filterBy/FilterButton.vue';
import {toValue} from '@/util/vue.js';
import {computed, ref, watch} from 'vue';

const props = defineProps({
  activeType: {
    type: String,
    default: ''
  },
  availableSources: {
    type: Array,
    default: () => []
  }
});

const filterInputValue = ref('');
const showHideFilterList = ref(false);
const selectedOption = ref(null);

// Compute the available sources based on the filter input value
const availableSearchedSources = computed(() => {
  if (!filterInputValue.value) return props.availableSources;

  return props.availableSources.filter((source) => {
    return source.title.toLowerCase().match(filterInputValue.value.toLowerCase()) ||
        source.title.toLowerCase().includes(filterInputValue.value.toLowerCase());
  });
});

const searchedSelectedOption = computed(() => {
  if (selectedOption.value !== null) {
    if (!filterInputValue.value) return toValue(selectedOption).value;

    return toValue(selectedOption).value.filter((source) => {
      return source.toLowerCase().match(filterInputValue.value.toLowerCase()) ||
          source.toLowerCase().includes(filterInputValue.value.toLowerCase());
    });
  } else {
    return [];
  }
});

// Reset the filter input value when the selected option changes
watch(selectedOption, () => {
  filterInputValue.value = '';
});
watch(showHideFilterList, () => {
  if (!showHideFilterList.value) {
    selectedOption.value = null;
  }
});
</script>

<style scoped>
.sticky-wrapper {
  position: sticky;
  top: 0; /* Adjust this value based on your needs */
  z-index: 1000; /* Ensure the sticky element stays above other content */
}
</style>
