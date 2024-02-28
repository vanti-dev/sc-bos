import {defineStore} from 'pinia';
import {computed, ref} from 'vue';

export const useNotificationFilterStore = defineStore('notificationFilter', () => {
  const severityLevels = {
    0: 'SEVERITY_UNSPECIFIED',
    9: 'INFO',
    13: 'WARN',
    17: 'ALERT',
    21: 'DANGER'
  };
  //
  // -------------------- Structure and Update Filter Query -------------------- //
  //
  // Active query is the reference to the current query
  // This can be directly modified by the user on import via filter changes
  const activeQuery = ref({
    acknowledged: undefined,
    createdNotBefore: undefined,
    createdNotAfter: undefined,
    severityNotAbove: undefined,
    severityNotBelow: undefined,
    floor: undefined,
    subsystem: undefined,
    source: undefined,
    resolved: false,
    resolvedNotBefore: undefined,
    resolvedNotAfter: undefined,
    zone: undefined
  });


  /**
   * Update the active query reference with the new value
   * This will only update the value if the key exists in the current context
   * Otherwise, it will be ignored
   * This is used to update the query on initial load - eg.: props are defined by the parent component
   *
   * @param {string} propName
   * @param {*} propValue
   * @return {void}
   */
  const updateActiveQuery = (propName, propValue) => {
    // Check if this key/value exists in the current context
    if (activeQuery.value.hasOwnProperty(propName)) {
      activeQuery.value[propName] = propValue; // Update the value
    }
  };

  /**
   * Reset the active query to the default state
   *
   * @return {void}
   */
  const resetActiveQuery = () => {
    Object.keys(activeQuery.value).forEach((key) => {
      if (key === 'resolved') {
        activeQuery.value[key] = false;
      } else {
        activeQuery.value[key] = undefined;
      }
    });
  };

  // -------------------- Available Sources -------------------- //
  /**
   * Reference to the available sources
   * This is used to populate the filter menu's main filter option list
   *
   * @type {import('vue').Ref<Array<*>>}
   */
  const availableSources = ref([]);

  /**
   * Reset the available sources to an empty array
   *
   * @return {void}
   */
  const resetAvailableSources = () => {
    availableSources.value = [];
  };

  // -------------------- Active Filters -------------------- //
  /**
   * Computed property that returns an array of the active filters
   *
   * @type {import('vue').ComputedRef<{value: *, key: *}[]>}
   */
  const activeFilters = computed(() => {
    return Object.entries(activeQuery.value)
        // filter out undefined or null values
        .filter(([key, value]) => value !== undefined && value !== null && key !== 'resolved')
        .map(([key, value]) => {
          if (key === 'acknowledged') {
            return {'key': key, 'value': value ? 'Yes' : 'No'};
          }

          if (key === 'severityNotAbove' || key === 'severityNotBelow') {
            return {'key': key, 'value': value};
          }

          return {'key': key, 'value': value};
        }); // map to an object
  });

  /**
   * Update the filter value in the active query
   *
   * @param {string} filter
   * @param {string|boolean|undefined} value
   * @return {void}
   */
  const updateFilter = (filter, value) => {
    if (filter === 'acknowledged') {
      activeQuery.value[filter] = value === 'Yes';
      return;
    }

    if (filter === 'severityNotAbove' || filter === 'severityNotBelow') {
      activeQuery.value[filter] = value;
      return;
    }

    activeQuery.value[filter] = value;
  };

  /**
   * Remove the filter from the active query
   *
   * @param {string} filter
   * @return {void}
   */
  const removeFilter = (filter) => {
    if (filter === 'resolved') {
      activeQuery.value[filter] = false;
      return;
    }

    activeQuery.value[filter] = undefined;
  };

  return {
    severityLevels,
    updateActiveQuery,
    resetActiveQuery,
    query: activeQuery.value,

    availableSources,
    resetAvailableSources,

    activeFilters,
    updateFilter,
    removeFilter
  };
});
