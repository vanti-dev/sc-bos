import {computed, reactive, ref} from 'vue';

/**
 * Returns attrs to v-bind to a VDataTableServer to display a collection.
 *
 * @param {Ref<number>} wantCount - the same ref used as part of UseCollectionOptions,
 *   will be updated based on table paging.
 * @param {UseCollectionResponse} collection - the response from useCollection
 * @return {Record<string, any>} attrs to v-bind to a VDataTableServer
 */
export function useDataTableCollection(wantCount, collection) {
  const fetchMoreItems = ({page, itemsPerPage}) => {
    wantCount.value = page * itemsPerPage;
  };
  const currentPage = ref(1);
  const itemsPerPage = ref(20);
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
    'onUpdate:items-per-page': (value) => itemsPerPage.value = value,
    'loading': collection.loadingNextPage
  });
}
