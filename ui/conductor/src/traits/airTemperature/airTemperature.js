import {closeResource, newActionTracker, newResourceValue} from '@/api/resource';
import {
  airTemperatureModeToString,
  pullAirTemperature,
  temperatureToString,
  updateAirTemperature
} from '@/api/sc/traits/air-temperature';
import {setRequestName, toQueryObject, watchResource} from '@/util/traits';
import {AirTemperature} from '@smart-core-os/sc-api-grpc-web/traits/air_temperature_pb';
import {computed, onScopeDispose, reactive, ref, toRefs, toValue} from 'vue';

/**
 * @typedef {
 *   import('@smart-core-os/sc-api-grpc-web/traits/air_temperature_pb').PullAirTemperatureRequest
 * } PullAirTemperatureRequest
 * @typedef {
 *   import('@smart-core-os/sc-api-grpc-web/traits/air_temperature_pb').PullAirTemperatureResponse
 * } PullAirTemperatureResponse
 * @typedef {
 *   import('@smart-core-os/sc-api-grpc-web/traits/air_temperature_pb').UpdateAirTemperatureRequest
 * } UpdateAirTemperatureRequest
 * @typedef {
 *  import('@smart-core-os/sc-api-grpc-web/traits/air_temperature_pb').AirTemperature
 * } AirTemperature
 * @typedef {import('vue').Ref} Ref
 * @typedef {import('vue').UnwrapNestedRefs} UnwrapNestedRefs
 * @typedef {import('vue').ToRefs} ToRefs
 * @typedef {import('vue').ComputedRef} ComputedRef
 * @typedef {import('@/api/resource').ResourceValue} ResourceValue
 * @typedef {import('@/api/resource').ActionTracker} ActionTracker
 */

/**
 * @param {MaybeRefOrGetter<string|PullAirTemperatureRequest.AsObject>} query
 * @param {MaybeRefOrGetter<boolean>=} paused
 * @return {ToRefs<UnwrapNestedRefs<ResourceValue<AirTemperature.AsObject, PullAirTemperatureResponse>>>}
 */
export function usePullAirTemperature(query, paused = false) {
  const resource = reactive(
      /** @type {ResourceValue<AirTemperature.AsObject, PullAirTemperatureResponse>} */
      newResourceValue()
  );
  onScopeDispose(() => closeResource(resource));

  const queryObject = computed(() => toQueryObject(query));

  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      (req) => {
        pullAirTemperature(req, resource);
        return resource;
      }
  );

  return toRefs(resource);
}

/**
 * @typedef UpdateAirTemperatureRequestLike
 * @type {number|Partial<AirTemperature.AsObject>|Partial<UpdateAirTemperatureRequest.AsObject>}
 */
/**
 * @param {MaybeRefOrGetter<string>=} name The name of the device to update.
 *   If not provided request objects must include a name.
 * @return {ToRefs<ActionTracker<AirTemperature.AsObject>> & {
 *   updateAirTemperature: (req: MaybeRefOrGetter<UpdateAirTemperatureRequestLike>) => Promise<AirTemperature.AsObject>
 * }}
 */
export function useUpdateAirTemperature(name) {
  const tracker = reactive(
      /** @type {ActionTracker<AirTemperature.AsObject>} */
      newActionTracker()
  );

  /**
   * @param {MaybeRefOrGetter<UpdateAirTemperatureRequestLike>} req
   * @return {UpdateAirTemperatureRequest.AsObject}
   */
  const toRequestObject = (req) => {
    req = toValue(req);
    if (typeof req === 'number') {
      req = {
        state: {temperatureSetPoint: {valueCelsius: /** @type {number} */ req}},
        updateMask: {pathsList: ['temperature_set_point']}
      };
    }
    if (!req.hasOwnProperty('state')) {
      req = {state: /** @type {AirTemperature.AsObject} */ req};
    }
    return setRequestName(req, name);
  };

  return {
    ...toRefs(tracker),
    updateAirTemperature: (req) => {
      return updateAirTemperature(toRequestObject(req), tracker);
    }
  };
}

/**
 * @param {MaybeRefOrGetter<AirTemperature.AsObject>} value
 * @return {{
 *   temp: ComputedRef<number>,
 *   setPoint: ComputedRef<number>,
 *   hasTemp: ComputedRef<boolean>,
 *   hasSetPoint: ComputedRef<boolean>,
 *   setPointStr: ComputedRef<string>,
 *   airTempData: ComputedRef<{}>,
 *   tempStr: ComputedRef<string>,
 *   tempRange: Ref<{low: number, high: number}>,
 *   tempProgress: ComputedRef<number>,
 *   humidity: ComputedRef<number>
 * }}
 */
export function useAirTemperature(value) {
  const _v = computed(() => toValue(value));

  const UNIT = '°C';
  const hasSetPoint = computed(() => {
    return _v.value?.temperatureSetPoint?.valueCelsius !== undefined;
  });
  const setPoint = computed(() => {
    return _v.value?.temperatureSetPoint?.valueCelsius;
  });
  const setPointStr = computed(() => {
    return setPoint.value?.toFixed(1) + UNIT;
  });

  const hasTemp = computed(() => {
    return _v.value?.ambientTemperature?.valueCelsius !== undefined;
  });
  const temp = computed(() => {
    return _v.value?.ambientTemperature?.valueCelsius;
  });
  const tempStr = computed(() => {
    const numStr = temp.value?.toFixed(1);
    if (hasSetPoint.value) {
      return numStr;
    } else {
      return numStr + UNIT;
    }
  });

  const tempRange = ref({
    low: 18.0,
    high: 24.0
  });
  const tempProgress = computed(() => {
    let val = temp.value ?? 0;
    if (val > 0) {
      val -= tempRange.value.low;
      val = val / (tempRange.value.high - tempRange.value.low);
    }
    return val * 100;
  });
  const airTempData = computed(() => {
    if (_v.value) {
      const data = {};
      Object.entries(_v.value).forEach(([key, value]) => {
        if (value !== undefined) {
          switch (key) {
            case 'mode':
              if (value !== AirTemperature.Mode.MODE_UNSPECIFIED) {
                data[key] = airTemperatureModeToString(value);
              }
              break;
            case 'ambientTemperature': {
              data['currentTemp'] = temperatureToString(value);
              break;
            }
            case 'temperatureSetPoint': {
              data['setPoint'] = temperatureToString(value);
              break;
            }
            case 'ambientHumidity':
              if (value !== 0) {
                data['humidity'] = value.toFixed(1) + '%';
              }
              break;
            case 'dewPoint': {
              data[key] = temperatureToString(value);
              break;
            }
            default: {
              data[key] = value;
            }
          }
        }
      });
      return data;
    }
    return {};
  });

  const humidity = computed(() => {
    return _v.value?.ambientHumidity;
  });

  return {
    hasSetPoint,
    setPoint,
    setPointStr,
    hasTemp,
    temp,
    tempStr,
    tempRange,
    tempProgress,
    airTempData,
    humidity
  };
}
