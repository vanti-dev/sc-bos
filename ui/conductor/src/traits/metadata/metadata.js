import {closeResource, newResourceValue} from '@/api/resource.js';
import {pullMetadata} from '@/api/sc/traits/metadata.js';
import {toQueryObject, watchResource} from '@/util/traits.js';
import {computed, onScopeDispose, reactive, toRefs, toValue} from 'vue';

/**
 * @typedef {
 *   import('@smart-core-os/sc-api-grpc-web/traits/metadata_pb').PullMetadataRequest
 * } PullMetadataRequest
 * @typedef {
 *   import('@smart-core-os/sc-api-grpc-web/traits/metadata_pb').PullMetadataResponse
 * } PullMetadataResponse
 * @typedef {
 *   import('@smart-core-os/sc-api-grpc-web/traits/metadata_pb').Metadata
 * } Metadata
 * @typedef {import('vue').Ref} Ref
 * @typedef {import('vue').UnwrapNestedRefs} UnwrapNestedRefs
 * @typedef {import('vue').ToRefs} ToRefs
 * @typedef {import('vue').ComputedRef} ComputedRef
 * @typedef {import('@/api/resource').ResourceValue} ResourceValue
 * @typedef {import('@/api/resource').ActionTracker} ActionTracker
 */

/**
 * @param {MaybeRefOrGetter<string|PullMetadataRequest.AsObject>} query
 * @param {MaybeRefOrGetter<boolean>=} paused
 * @return {ToRefs<ResourceValue<Metadata.AsObject, PullMetadataResponse>>}
 */
export function usePullMetadata(query, paused = false) {
  const resource = reactive(
      /** @type {ResourceValue<Metadata.AsObject, PullMetadataResponse>} */
      newResourceValue()
  );
  onScopeDispose(() => closeResource(resource));

  const queryObject = computed(() => toQueryObject(query));

  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      (req) => {
        pullMetadata(req, resource);
        return () => closeResource(resource);
      }
  );

  return toRefs(resource);
}
