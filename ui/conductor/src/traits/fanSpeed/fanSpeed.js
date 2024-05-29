import {closeResource, newResourceValue} from '@/api/resource.js';
import {pullFanSpeed} from '@/api/sc/traits/fan-speed.js';
import {toQueryObject, watchResource} from '@/util/traits.js';
import {toValue} from '@/util/vue.js';
import {computed, onScopeDispose, reactive, toRefs} from 'vue';

/**
 * @param {MaybeRefOrGetter<string|PullFanSpeedRequest.AsObject>} query
 * @param {MaybeRefOrGetter<boolean>=} paused
 * @return {ToRefs<ResourceValue<FanSpeed.AsObject, PullFanSpeedResponse>>}
 */
export function usePullFanSpeed(query, paused = false) {
  const resource = reactive(
      /** @type {ResourceValue<FanSpeed.AsObject, PullFanSpeedResponse>} */
      newResourceValue());
  onScopeDispose(() => closeResource(resource));
  const queryObject = computed(() => toQueryObject(query));
  watchResource(() => toValue(queryObject), () => toValue(paused), req => {
    pullFanSpeed(req, resource);
    return resource;
  });
  return toRefs(resource);
}
