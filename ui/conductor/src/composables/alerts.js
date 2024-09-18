import {timestampToDate} from '@/api/convpb.js';
import {closeResource, newResourceValue} from '@/api/resource.js';
import {listAlerts, pullAlertMetadata, pullAlerts} from '@/api/ui/alerts.js';
import useCollection from '@/composables/collection.js';
import {cmpDesc} from '@/util/date.js';
import {toQueryObject, watchResource} from '@/util/traits.js';
import {computed, onScopeDispose, reactive, toRefs, toValue} from 'vue';

/**
 * @param {MaybeRefOrGetter<Partial<ListAlertsRequest.AsObject & PullAlertsRequest.AsObject>>} request
 * @param {MaybeRefOrGetter<Partial<UseCollectionOptions<Alert.AsObject>>>?} options
 * @return {UseCollectionResponse<Alert.AsObject>}
 */
export function useAlertsCollection(request, options) {
  const normOpts = computed(() => {
    return {
      cmp: (a, b) => cmpDesc(timestampToDate(a.createTime), timestampToDate(b.createTime)),
      ...toValue(options)
    };
  });
  const client = {
    async listFn(req, tracker) {
      const res = await listAlerts(req, tracker);
      return {
        items: res.alertsList,
        nextPageToken: res.nextPageToken,
        totalSize: res.totalSize
      };
    },
    pullFn(req, resource) {
      pullAlerts(req, resource);
    }
  };
  return useCollection(request, client, normOpts);
}

/**
 * @param {MaybeRefOrGetter<string|Partial<PullAlertMetadataRequest.AsObject>>} request
 * @param {MaybeRefOrGetter<{paused?: boolean}>?} options
 * @return {ToRefs<ResourceValue<AlertMetadata.AsObject, any>>}
 */
export function usePullAlertMetadata(request, options) {
  const resource = reactive(
      /** @type {ResourceValue<AlertMetadata.AsObject, PullAlertMetadataResponse>} */
      newResourceValue()
  );
  onScopeDispose(() => closeResource(resource));

  const queryObject = computed(() => toQueryObject(request));

  watchResource(
      () => toValue(queryObject),
      () => toValue(options)?.paused ?? false,
      (req) => {
        pullAlertMetadata(req, resource);
        return () => closeResource(resource);
      }
  );

  return toRefs(resource);
}
