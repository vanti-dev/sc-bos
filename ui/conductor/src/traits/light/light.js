import {closeResource, newActionTracker, newResourceValue} from '@/api/resource';
import {describeBrightness, pullBrightness, updateBrightness} from '@/api/sc/traits/light';
import {setRequestName, toQueryObject, watchResource} from '@/util/traits.js';
import {toValue} from '@/util/vue';
import {computed, onUnmounted, reactive, toRefs} from 'vue';

/**
 * @typedef {import('@smart-core-os/sc-api-grpc-web/traits/light_pb').PullBrightnessRequest} PullBrightnessRequest
 * @typedef {import('@smart-core-os/sc-api-grpc-web/traits/light_pb').PullBrightnessResponse} PullBrightnessResponse
 * @typedef {import('@smart-core-os/sc-api-grpc-web/traits/light_pb').UpdateBrightnessRequest} UpdateBrightnessRequest
 * @typedef {
 *   import('@smart-core-os/sc-api-grpc-web/traits/light_pb').DescribeBrightnessRequest
 * } DescribeBrightnessRequest
 * @typedef {import('@smart-core-os/sc-api-grpc-web/traits/light_pb').Brightness} Brightness
 * @typedef {import('@smart-core-os/sc-api-grpc-web/traits/light_pb').BrightnessSupport} BrightnessSupport
 * @typedef {import('@smart-core-os/sc-api-grpc-web/traits/light_pb').LightPreset} LightPreset
 * @typedef {import('vue').Ref} Ref
 * @typedef {import('vue').ToRefs} ToRefs
 * @typedef {import('vue').ComputedRef} ComputedRef
 * @typedef {import('@/api/resource').ResourceValue} ResourceValue
 * @typedef {import('@/api/resource').ActionTracker} ActionTracker
 */

/**
 * @param {MaybeRefOrGetter<string|PullBrightnessRequest.AsObject>} query
 * @param {MaybeRefOrGetter<boolean>=} paused
 * @return {ToRefs<ResourceValue<Brightness.AsObject, PullBrightnessResponse>>}
 */
export function usePullBrightness(query, paused = false) {
  const resource = reactive(
      /** @type {ResourceValue<Brightness.AsObject, PullBrightnessResponse>} */
      newResourceValue()
  );
  onUnmounted(() => closeResource(resource));

  const queryObject = computed(() => toQueryObject(query));

  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      (req) => {
        pullBrightness(req, resource);
        return resource;
      }
  );

  return toRefs(resource);
}

/**
 * @param {MaybeRefOrGetter<string|DescribeBrightnessRequest.AsObject>} query
 * @param {MaybeRefOrGetter<number>=} pollInterval
 * @return {ToRefs<ActionTracker<BrightnessSupport.AsObject>>}
 */
export function useDescribeBrightness(query, pollInterval = 0) {
  const tracker = reactive(
      /** @type {ActionTracker<BrightnessSupport.AsObject>} */
      newActionTracker()
  );

  const queryObject = computed(() => toQueryObject(query));

  watchResource(
      () => toValue(queryObject),
      () => false,
      (req) => {
        describeBrightness(req, tracker)
            .catch(() => {}); // errors are handled by the tracker
        return tracker;
      }
  );

  return toRefs(tracker);
}

/**
 * @typedef UpdateBrightnessRequestLike
 * @type {number|Partial<Brightness.AsObject>|Partial<UpdateBrightnessRequest.AsObject>}
 */
/**
 * @param {MaybeRefOrGetter<string>} name
 * @return {ToRefs<ActionTracker<Brightness.AsObject>> & {
 *  updateBrightness: (req: MaybeRefOrGetter<UpdateBrightnessRequestLike>) => Promise<Brightness.AsObject>
 * }}
 */
export function useUpdateBrightness(name) {
  const tracker = reactive(
      /** @type {ActionTracker<Brightness.AsObject>} */
      newActionTracker()
  );

  const toRequestObject = (req) => {
    req = toValue(req);
    if (typeof req === 'number') {
      req = {
        brightness: {levelPercent: req},
        updateMask: {pathsList: ['level_percent']}
      };
    }
    if (!req.hasOwnProperty('brightness')) {
      req = {brightness: req};
    }
    return setRequestName(req, name);
  };

  return {
    ...toRefs(tracker),
    updateBrightness: (req) => {
      return updateBrightness(toRequestObject(req), tracker);
    }
  };
}

/**
 * @param {MaybeRefOrGetter<Brightness.AsObject>} value
 * @param {MaybeRefOrGetter<BrightnessSupport.AsObject|null>=} support
 * @return {{
 *   level: ComputedRef<number>,
 *   levelStr: ComputedRef<string>,
 *   icon: ComputedRef<string>,
 *   presets: ComputedRef<Array<LightPreset.AsObject>>
 * }}
 */
export function useBrightness(value, support = null) {
  const _v = computed(() => toValue(value));
  const _s = computed(() => toValue(support));

  /** @type {ComputedRef<number>} */
  const level = computed(() => _v.value?.levelPercent || 0);

  /** @type {ComputedRef<string>} */
  const levelStr = computed(() => {
    if (level.value === 0) {
      return 'Off';
    } else if (level.value === 100) {
      return 'Max';
    } else if (level.value > 0 && level.value < 100) {
      return `${level.value.toFixed(0)}%`;
    }

    return '';
  });

  /** @type {ComputedRef<string>} */
  const icon = computed(() => {
    if (level.value === 0) {
      return 'mdi-lightbulb-outline';
    } else if (level.value > 0) {
      return 'mdi-lightbulb-on';
    } else {
      return '';
    }
  });


  /** @type {ComputedRef<Array<LightPreset.AsObject>>} */
  const presets = computed(() => {
    return _s.value ?.presetsList ?? [];
  });

  return {
    level,
    levelStr,
    icon,
    presets
  };
}
