import {useNotificationFilterStore} from '@/routes/ops/notifications/useNotificationFilterStore.js';
import {toValue} from '@/util/vue.js';
import {computed} from 'vue';

/**
 * @typedef {import('vue').Ref<
 *  import('@/routes/ops/notifications/useNotificationFilterStore').NotificationFilterStore>
 * } NotificationFilterStore
 */

/**
 *
 * @param {MaybeRefOrGetter<string>} type
 * @return {{
 *  filterStore: import('vue').ComputedRef<NotificationFilterStore|null>
 * }}
 */
export default function(type) {
  const notificationStore = useNotificationFilterStore();

  /**
   * Computed property that returns the appropriate filter store based on the type of filter required
   *
   * @type {import('vue').ComputedRef<NotificationFilterStore|null>}
   */
  const filterStore = computed(() => {
    let store = null;

    if (toValue(type) === 'notification') {
      store = notificationStore;
    }

    return toValue(store);
  });

  // ----------------------------- Filters & Sources ----------------------------- //
  /**
   * Computed property that returns an array of the available sources
   *
   * @type {import('vue').ComputedRef<Array<*>>}
   */
  const availableSources = computed(() => {
    const store = toValue(filterStore);

    if (store) {
      return store.availableSources.sort((a, b) => a.title.localeCompare(b.title));
    }

    return [];
  });

  /**
   * Computed property that returns an array of the active filters
   *
   * @type {import('vue').ComputedRef<*|[]>}
   */
  const activeFilters = computed(() => {
    const store = toValue(filterStore);

    if (store) {
      return store.activeFilters.sort((a, b) => a.key.localeCompare(b.key));
    }

    return [];
  });

  // ----------------------------- Filter Modifiers ----------------------------- //
  /**
   * Update the filter value in the active query
   *
   * @param {string} filter
   * @param {string|boolean|undefined} value
   * @return {void}
   */
  const updateFilter = (filter, value) => {
    const store = toValue(filterStore);

    if (store) {
      if (['severityNotAbove', 'severityNotBelow'].includes(filter)) {
        store.updateFilter(filter, value);
        return;
      }

      store.updateFilter(filter.toLowerCase(), value);
    }
  };

  /**
   * Remove the filter from the active query
   *
   * @param {string} filter
   * @return {void}
   */
  const removeFilter = (filter) => {
    const store = toValue(filterStore);

    if (store) {
      if (['severityNotAbove', 'severityNotBelow'].includes(filter)) {
        store.removeFilter(filter);
        return;
      }

      store.removeFilter(filter.toLowerCase());
    }
  };

  return {
    filterStore,
    availableSources,
    activeFilters,

    updateFilter,
    removeFilter
  };
}
