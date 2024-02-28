<template>
  <v-row class="mr-0 d-flex flex-row align-center">
    <v-sheet v-if="props.activeFilters.length" class="px-1" elevation="0" max-width="575">
      <v-chip-group>
        <v-tooltip v-for="filter in props.activeFilters" :key="filter.key" bottom>
          <template #activator="{ on }">
            <v-chip
                class="text-capitalize text-truncate non-interactive-chip"
                :close="!disableRemove(filter.key)"
                close-label="Remove filter"
                :disabled="disableRemove(filter.key)"
                label
                :ripple="false"
                small
                style="max-width: 200px;"
                v-on="on"
                @click:close="disableRemove(filter.key) ? null : removeFilter(filter.key)">
              {{ formatObjectValues(filter).key }}: {{ formatObjectValues(filter).value }}
            </v-chip>
          </template>
          <span>Remove filter</span>
        </v-tooltip>
      </v-chip-group>
    </v-sheet>
    <v-subheader v-else>No filters applied.</v-subheader>
  </v-row>
</template>

<script setup>
import useFilterBy from '@/components/filterBy/useFilterBy.js';
import {useNotificationFilterStore} from '@/routes/ops/notifications/useNotificationFilterStore.js';
import {camelToSentence} from '@/util/string.js';

const props = defineProps({
  activeType: {
    type: String,
    default: ''
  },
  activeFilters: {
    type: Array,
    default: () => []
  },
  availableSources: {
    type: Array,
    default: () => []
  }
});

const notificationFilterStore = useNotificationFilterStore();
const {removeFilter} = useFilterBy(() => props.activeType);

/**
 * Disable remove action if filter is set by default and source not available
 *
 * @param {string} key
 * @return {boolean}
 */
const disableRemove = (key) => {
  if (key === 'severityNotBelow' || key === 'severityNotAbove') {
    return false;
  }

  // Check if the key exists in availableSources by comparing the title
  const source = props.availableSources.find(source => source.title.toLowerCase() === key.toLowerCase());

  return !source;
};


const formatObjectValues = (obj) => {
  let formattedObj;

  if (obj.key === 'severityNotBelow' || obj.key === 'severityNotAbove') {
    formattedObj = {
      key: camelToSentence(obj.key),
      value: notificationFilterStore.severityLevels[obj.value]
    };
  } else {
    formattedObj = {
      key: camelToSentence(obj.key),
      value: obj.value
    };
  }

  return formattedObj;
};
</script>

<style lang="scss" scoped>
.non-interactive-chip {
  pointer-events: none; /* Disable pointer events on the chip as they are not functioning */
}

.non-interactive-chip ::v-deep.v-chip:not(.v-chip--disabled) .v-chip__close {
  pointer-events: all; /* Re-enable pointer events on the close button so we can close the chip/remove filter */

  &:hover {
    color: var(--v-error-base); /* Change the close button color on hover for visual effect */
  }
}
</style>
