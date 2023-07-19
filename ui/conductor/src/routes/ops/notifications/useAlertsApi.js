/** @typedef {import('@sc-bos/ui-gen/proto/alerts_pb').Alert} Alert */
/** @typedef {import('@sc-bos/ui-gen/proto/alerts_pb').ListAlertsRequest} ListAlertsRequest */
/** @typedef {import('@sc-bos/ui-gen/proto/alerts_pb').ListAlertsResponse} ListAlertsResponse */
/** @typedef {import('@sc-bos/ui-gen/proto/alerts_pb').PullAlertsRequest} PullAlertsRequest */
/** @typedef {import('@sc-bos/ui-gen/proto/alerts_pb').PullAlertsResponse} PullAlertsResponse */

import {timestampToDate} from '@/api/convpb';
import {closeResource, newActionTracker, newResourceCollection} from '@/api/resource';
import {listAlerts, pullAlerts} from '@/api/ui/alerts';
import {toValue} from '@/util/vue';
import {computed, onMounted, onUnmounted, reactive, ref, watch} from 'vue';

/**
 * @param {MaybeRefOrGetter<string>} name
 * @param {MaybeRefOrGetter<Alert.Query.AsObject>} query
 * @return {{}}
 */
export default function(name, query) {
  // Alerts we've fetched via ListAlerts, in the order based on comparator, transformed
  const listedItems = ref(/** @type {JSAlert[]} */ []);
  // Items in pullResource, sorted by createTime, and transformed.
  const pulledItems = computed(() => {
    return Object.values(pullResource.value ?? {})
        .map(i => transform(i))
        .sort(comparator);
  });
  // Both items created via pull and those listed via paging, in order with dupes removed.
  const allItems = computed(() => {
    const listed = listedItems.value;
    const pulled = pulledItems.value;
    const items = /** @type {JSAlert[]} */ [];

    let li = 0;
    let pi = 0;
    while (li < listed.length && pi < pulled.length) {
      const cmp = comparator(listed[li], pulled[pi]);
      if (cmp < 0) {
        items.push(listed[li++]);
      } else if (cmp === 0) {
        items.push(pulled[pi++]);
        li++;
      } else {
        items.push(pulled[pi++]);
      }
    }

    if (pi < pulled.length) {
      items.push(...pulled.slice(pi));
    } else {
      items.push(...listed.slice(li));
    }

    return items;
  });

  // how many items have we fetched from the server?
  const listedItemCount = computed(() => listedItems.value.length);
  // how many items are we hoping we'll get from the server?
  const targetItemCount = ref(0);

  // pagination state
  const pageSize = ref(100);
  const pageIndex = ref(0);
  watch([pageSize, pageIndex], () => {
    // could optimise this to not require items that we've pulled into the tip.
    targetItemCount.value = pageSize.value * (pageIndex.value + 1);
  }, {immediate: true});

  // Items on the page identified by pageSize and pageIndex
  const pageItems = computed(() => {
    const pageStart = pageSize.value * pageIndex.value;
    return allItems.value.slice(pageStart, pageStart + pageSize.value);
  });

  // tracks our pull request
  const pullResource = reactive(
      /** @type {ResourceCollection<Alert.AsObject, PullAlertsResponse>} */
      newResourceCollection());
  // tracks each fetch of a new page, resource value may be outdated
  const listPageTracker = reactive(
      /** @type {ActionTracker<ListAlertsResponse.AsObject>} */
      newActionTracker());
  const nextPageToken = ref('');

  // queries we run against the server
  const listQuery = computed(() => {
    return /** @type {ListAlertsRequest.AsObject} */ {
      name: toValue(name),
      query: toValue(query)
    };
  });
  const pullQuery = computed(() => {
    return /** @type {PullAlertsRequest.AsObject} */ {name: toValue(name), query: toValue(query)};
  });

  const mounted = ref(false);
  onMounted(() => mounted.value = true);
  onUnmounted(() => mounted.value = false);

  const hasFetchedAnyPages = ref(false);
  const shouldFetchMorePages = computed(() => {
    if (!mounted.value) return false;
    // do we want more alerts, and do we think there are more alerts?
    return listedItemCount.value < targetItemCount.value &&
        (!hasFetchedAnyPages.value || Boolean(nextPageToken.value));
  });

  // Debug property that keeps track of past ListAlerts requests we've made.
  const pastListRequests = ref(/** @type {ListAlertsRequest.AsObject[]} */ []);
  const recordListRequest = req => {
    pastListRequests.value.push(req);
    if (pastListRequests.value.length > 5) pastListRequests.value.shift();
  };

  // either false, or a number indicating which version of listQuery we're running
  const fetchingPages = ref(/** @type {boolean|number} */ false);
  // This is set if fetchMore is called with a new query version while a fetch is ongoing.
  const nextQuery = ref(
      /** @type {{query: ListAlertsRequest.AsObject, version: boolean} | null} */
      null);
  // Fetches more pages from the server to meet the targetItemCount.
  // Will loop fetching more pages until the alertCount is >= targetItemCount.
  const fetchMore = async (query, version) => {
    if (fetchingPages.value === version) return; // we're already fetching this query version
    if (fetchingPages.value !== false) {
      // we've been asked to fetch a page but we're already fetching a page, but with an old query
      // need to tidy things up and start again.
      nextQuery.value = {query, version};
      return; // nextQuery will be run when the existing fetch completes
    }

    try {
      query = {...query, pageSize: pageSize.value}; // clone so we don't mutate the original
      while (shouldFetchMorePages.value) {
        fetchingPages.value = version;
        query.pageToken = nextPageToken.value;
        recordListRequest({...query}); // for debugging
        const page = await listAlerts(query, listPageTracker);

        if (nextQuery.value) {
          // we were asked to fetch a new query while we were fetching a page, so we need to start again.
          query = {...nextQuery.value.query};
          version = nextQuery.value.version;
          nextQuery.value = null;
          closeResource(listPageTracker);
          continue;
        }

        // are the results still useful, i.e. is the query we started with still valid?
        // success case, we fetched a page and nothing updated while we waited.
        listedItems.value.push(...page.alertsList.map(a => transform(a)));
        nextPageToken.value = page.nextPageToken;
        hasFetchedAnyPages.value = true;
      }
    } finally {
      listedItems.value.sort(comparator); // make sure the sort order is consistent.
      fetchingPages.value = false;
    }
  };

  watch(
      [pullQuery],
      () => {
        closeResource(pullResource);
        const request = pullQuery.value;
        if (request) {
          pullAlerts(request, pullResource);
        }
      },
      {immediate: true, deep: true}
  );
  const queryVersionCounter = ref(0);
  watch(listQuery, () => {
    queryVersionCounter.value++;
    // tidy up state, if the query has changed then these are no longer valid.
    hasFetchedAnyPages.value = false;
    nextPageToken.value = '';
    listedItems.value = [];
  }, {deep: true});
  watch(
      [shouldFetchMorePages, queryVersionCounter],
      ([_, v]) => {
        fetchMore(listQuery.value, v)
            // errors are tracked by listPageTracker
            .catch(() => {
            });
      },
      {immediate: true}
  );

  const loading = computed(() => {
    return fetchingPages.value !== false || pullResource.loading;
  });

  return {
    listedItems,
    pulledItems,
    pageItems,
    allItems,
    loading,

    listedItemCount,
    // write to this to tell us how many items you want to read
    targetItemCount,

    // write these to adjust page related settings
    pageSize,
    pageIndex,

    // used for troubleshooting
    pullQuery,
    pullResource,
    listQuery,
    listPageTracker,
    nextPageToken,
    nextQuery,
    fetchingPages,
    shouldFetchMorePages,
    pastListRequests,
    queryVersionCounter
  };
}

/**
 * @typedef {Alert.AsObject & {createTime: Date, resolveTime: Date}} JSAlert
 */

/**
 * @param {Alert.AsObject} alert
 * @return {JSAlert}
 */
const transform = (alert) => {
  alert.createTime = timestampToDate(alert.createTime);
  alert.resolveTime = timestampToDate(alert.resolveTime);
  if (alert.acknowledgement) {
    alert.acknowledgement.acknowledgeTime = timestampToDate(alert.acknowledgement.acknowledgeTime);
  }
  return alert;
};

// createTime descending, with ties broken by id descending
const comparator = (a, b) => {
  const aTime = a.createTime.getTime();
  const bTime = b.createTime.getTime();
  if (aTime === bTime) return b.id.localeCompare(a.id);
  return bTime - aTime;
};
