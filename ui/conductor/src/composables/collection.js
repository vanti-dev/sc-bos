import {closeResource, newActionTracker, newResourceCollection} from '@/api/resource.js';
import {cap} from '@/util/number.js';
import {watchResource} from '@/util/traits.js';
import deepEqual from 'fast-deep-equal';
import {computed, reactive, ref, toValue, watch} from 'vue';

/**
 * @typedef {Object} ListResponse
 * @template T
 * @property {T[]} items
 * @property {string=} nextPageToken
 * @property {number=} totalSize
 */
/**
 * @typedef {Object} PageRequest
 * @property {number=} pageSize
 * @property {string=} pageToken
 */
/**
 * @typedef {Object} PullChange
 * @template T
 * @property {T} oldValue
 * @property {T} newValue
 */
/**
 * @typedef {Object} UseCollectionOptions
 * @template T
 * @property {number=} wantCount - how many items to fetch from the server, -1 for all
 * @property {number=} pageSize - how many items to fetch per request, defaults to cap(missing, 10, 500)
 * @property {boolean=} paused - suspend requests
 * @property {(item: T) => string=} idFn - a function to get the id of an item, defaults to item.id or item.name
 * @property {(item: T) => T} transform - a function that transforms items before they are added to the list.
 *   Useful for converting proto Timestamps.
 * @property {(a: T, b: T) => number} cmp - a function that compares two items. Any new item we receive from the
 *   subscription will not be added unless this is specified, as we don't know where to put it in the list.
 */
/**
 * @typedef {Object} UseCollectionResponse
 * @template T
 * @property {Ref<T[]>} items
 * @property {Ref<number>} totalItems
 * @property {Ref<boolean>} hasServerTotalItems
 * @property {Ref<boolean>} hasMorePages
 * @property {Ref<boolean>} loading
 * @property {Ref<boolean>} loadingNextPage
 * @property {Ref<ResourceError[]>} errors
 */

/**
 * Executes a list query and keeps the results up to date based on a pull.
 *
 * @template T, R
 * @param {MaybeRefOrGetter<R>} request - the request to use for the collection. Passed to both listFn and pullFn,
 *   for listFn it may have pageSize or pageToken set.
 * @param {{
 *   listFn: (req: R & PageRequest, tracker: ActionTracker<any>) => Promise<ListResponse<T>>,
 *   pullFn?: (req: R, resource: ResourceCollection<T, any>) => void
 * }} client - functions for listing and subscribing to changes based on a request
 * @param {MaybeRefOrGetter<UseCollectionOptions<T>>} [options] - optional options
 * @return {UseCollectionResponse<T>}
 */
export default function useCollection(request, client, options) {
  // The final items of the collection, combining listed and pulled items.
  const items = ref(/** @type {T[]} */ []);

  // related to calls to client.listFn
  const listTracker = reactive(/** @type {ActionTracker<T>} */ newActionTracker());
  const lastListResponse = ref(/** @type {ListResponse<T>} */ null);

  // related to calls to client.pullFn
  const pullResource = reactive(/** @type {ResourceCollection<T, any>} */ newResourceCollection());
  // changes we pulled while client.listFn was fetching, they will be applied once listFn is done
  const pullChanges = ref(/** @type {PullChange<T>[]} */ []);
  watch(() => pullResource.lastResponse, (r) => {
    // most pull responses look like {changesList: [{oldValue, newValue}]}
    if (r && typeof r.toObject === 'function' && typeof r.getChangesList === 'function') {
      pullChanges.value.push(...r.getChangesList().map(change => change.toObject()));
      if (!listTracker.loading) processChanges();
    }
  }, {flush: 'sync'});


  const targetListCount = computed(() => toValue(options)?.wantCount ?? 20);
  // Are there (likely) more pages available on the server?
  // If we've never asked, or the server says there are more pages, then we return true.
  const hasMorePages = computed(() => !lastListResponse.value || !!lastListResponse.value.nextPageToken);
  const shouldFetch = computed(() => {
    if (toValue(options)?.paused ?? false) return false; // don't fetch if paused
    if (listTracker.loading) return false; // don't fetch if already fetching
    if (!hasMorePages.value) return false; // don't fetch if there are no more pages
    // otherwise, fetch if we haven't fetched enough items
    return targetListCount.value === -1 || items.value.length < targetListCount.value;
  });
  // A guess at how many total items there are, either from the server or calculated locally based on fetched items.
  const totalItems = computed(() => lastListResponse.value?.totalSize ?? items.value.length);
  // Is totalItems a value returned by the server or calculated locally.
  const hasServerTotalItems = computed(() => Boolean(lastListResponse.value?.totalSize));

  const loading = computed(() => listTracker.loading || pullResource.loading);
  const loadingNextPage = computed(() => listTracker.loading);
  const errors = computed(() => [listTracker.error, pullResource.streamError]
      .filter(e => e));

  // data fetching
  watchResource(
      () => toValue(request),
      () => toValue(options)?.paused,
      (req) => {
        client.pullFn(req, pullResource);
        return () => closeResource(pullResource);
      }
  );
  const shouldFetchWatcherRunning = ref(false);
  watch(shouldFetch, async () => {
    if (shouldFetchWatcherRunning.value) return;
    shouldFetchWatcherRunning.value = true;
    try {
      while (shouldFetch.value) {
        await fetchNextPage();
        processChanges();
      }
    } catch (e) {
      // todo: add options to not log the error because the caller is handling it
      console.warn(e);
    } finally {
      shouldFetchWatcherRunning.value = false;
    }
  }, {immediate: true});
  watch(() => toValue(request), (o, n) => {
    if (deepEqual(o, n)) return; // no change
    items.value = [];
    lastListResponse.value = null;
    // the change in request will also cause the pull watcher to trigger,
    // but there's no way in there to know if we were paused or the request changed.
    pullChanges.value = [];
  }, {deep: true});

  /**
   * Calls client.listFn to fetch the next page of items.
   * Sets lastListResponse and updates listTracker.
   *
   * @return {Promise<void>}
   */
  async function fetchNextPage() {
    const _request = {...toValue(request)}; // clone so we can modify it
    if (lastListResponse.value) _request.pageToken = lastListResponse.value.nextPageToken;
    _request.pageSize = _request.pageSize ??
        toValue(options)?.pageSize ??
        cap(targetListCount.value - items.value.length, 10, 500);

    const pageResponse = await client.listFn(_request, listTracker);
    lastListResponse.value = pageResponse;

    const transform = toValue(options)?.transform ?? (v => v);
    items.value.push(...pageResponse.items.map(transform));
  }

  /**
   * Processes changes we received from a pull against those we received from a list.
   * pullChanges will be empty after this call.
   *
   * If options.cmp is not defined, items created via pull will not be added to the list.
   */
  function processChanges() {
    const _items = items.value;
    const _changes = pullChanges.value;
    if (!_changes.length) return; // no changes, nothing to do
    const transform = toValue(options)?.transform ?? (v => v);

    // optimise the lookup of items by id if we're going to do it a bunch
    let getIndex = (id) => _items.findIndex(v => getId(v) === id);
    if (_changes.length > 10) {
      // compute an index for the items
      const index = new Map();
      for (let i = 0; i < _items.length; i++) {
        index.set(getId(_items[i]), i);
      }
      getIndex = (id) => index.get(id);
    }

    // delay mutating the items until later so the indexes remain accurate
    const toDeleteIndexes = [];
    const toAddItems = [];

    for (const change of pullChanges.value) {
      const index = getIndex(getId(change.newValue ?? change.oldValue));
      if (change.oldValue && !change.newValue) {
        // deletion
        if (index !== -1) toDeleteIndexes.push(index);
      } else if (change.newValue) {
        if (index !== -1) {
          // update
          _items[index] = transform(change.newValue);
        } else {
          // add
          toAddItems.push(transform(change.newValue));
        }
      }
    }

    // process deletes first
    toDeleteIndexes.sort((a, b) => b - a); // sort in reverse order
    for (const i of toDeleteIndexes) {
      _items.splice(i, 1);
    }

    // only insert if options.cmp is defined, then insert in the correct place
    const cmp = toValue(options)?.cmp;
    if (cmp) {
      for (const item of toAddItems) {
        let i = 0;
        while (i < _items.length && cmp(item, _items[i]) > 0) i++;
        _items.splice(i, 0, item);
      }
    } else {
      // don't actually add the items as we'd have no way to know where to put them
    }

    items.value = _items;
    pullChanges.value = [];
  }

  return {
    items,
    totalItems,
    hasServerTotalItems,
    hasMorePages,

    loading,
    loadingNextPage,
    errors,

    _listTracker: listTracker,
    _pullResource: pullResource
  };
}

/**
 * Get the id for the value. Uses idFn or .id or .name or returns null.
 *
 * @template T
 * @param {T} v
 * @param {(T) => string=} idFn
 * @return {string|null}
 */
function getId(v, idFn) {
  if (idFn) return idFn(v);
  return v?.id ?? v?.name ?? null;
}
