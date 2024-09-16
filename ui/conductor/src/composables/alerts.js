import {closeResource, newResourceValue} from '@/api/resource.js';
import {pullAlertMetadata} from '@/api/ui/alerts.js';
import {toQueryObject, watchResource} from '@/util/traits.js';
import {computed, onScopeDispose, reactive, toRefs, toValue} from 'vue';

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
