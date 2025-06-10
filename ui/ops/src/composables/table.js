import {useLocalProp} from '@/util/vue.js';
import binarySearch from 'binary-search';
import {computed, reactive, ref, toValue, watch} from 'vue';

/**
 * @typedef {Object} useDataTableCollectionOptions
 * @property {number} [itemsPerPage=20] - the number of items to display per page initially
 */

/**
 * Returns attrs to v-bind to a VDataTableServer to display a collection.
 *
 * @param {Ref<number>} wantCount - the same ref used as part of UseCollectionOptions,
 *   will be updated based on table paging.
 * @param {UseCollectionResponse} collection - the response from useCollection
 * @param {import('vue').MaybeRefOrGetter<useDataTableCollectionOptions>} [options] - options for the collection
 * @return {Record<string, any>} attrs to v-bind to a VDataTableServer
 */
export function useDataTableCollection(wantCount, collection, options = null) {
  const fetchMoreItems = ({page, itemsPerPage}) => {
    wantCount.value = page * itemsPerPage;
  };
  const currentPage = ref(1);
  const itemsPerPage = useLocalProp(computed(() => toValue(options)?.itemsPerPage ?? 20));
  const itemsPerPageOptions = ref([...defaultItemsPerPageOptions]);
  watch(itemsPerPage, (itemsPerPage) => {
    const idx = binarySearch(itemsPerPageOptions.value, itemsPerPage, (a, b) => a.value - b);
    if (idx < 0) {
      // not found, add it to the list
      itemsPerPageOptions.value.splice(~idx, 0, {title: String(itemsPerPage), value: itemsPerPage});
    }
    // note: we aren't removing old items because we don't want to forget them if the user selects another from the options.
  }, {immediate: true});

  const pagedItems = computed(() => {
    const start = (currentPage.value - 1) * itemsPerPage.value;
    const end = currentPage.value * itemsPerPage.value;
    return collection.items.value.slice(start, end);
  });

  // as we intend the result to be used like v-bind="useDataTableCollection(...)",
  // we have to return a reactive object as v-bind doesn't work with plain objects containing refs.
  return reactive({
    'items': pagedItems,
    'onUpdate:options': fetchMoreItems,
    'items-length': collection.totalItems,
    'page': currentPage,
    'onUpdate:page': (value) => currentPage.value = value,
    'items-per-page': itemsPerPage,
    'items-per-page-options': itemsPerPageOptions,
    'onUpdate:items-per-page': (value) => itemsPerPage.value = value,
    'loading': collection.loadingNextPage
  });
}

export const defaultItemsPerPageOptions = [
  {title: '20', value: 20},
  {title: '50', value: 50},
  {title: '100', value: 100},
]