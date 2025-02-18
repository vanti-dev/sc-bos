import {closeResource, newResourceValue} from '@/api/resource.js';
import {pullAirQualitySensor} from '@/api/sc/traits/air-quality-sensor';
import {toQueryObject, watchResource} from '@/util/traits.js';
import {AirQuality} from '@smart-core-os/sc-api-grpc-web/traits/air_quality_sensor_pb';
import {computed, onScopeDispose, reactive, toRefs, toValue} from 'vue';

/**
 * @typedef {
 *   import('@smart-core-os/sc-api-grpc-web/traits/air_quality_sensor_pb').PullAirQualityRequest
 * } PullAirQualityRequest
 * @typedef {
 *   import('@smart-core-os/sc-api-grpc-web/traits/air_quality_sensor_pb').PullAirQualityResponse
 * } PullAirQualityResponse
 * @typedef {import('vue').ComputedRef} ComputedRef
 * @typedef {import('@/api/resource').ResourceError} ResourceError
 * @typedef {import('@/api/resource').ResourceValue} ResourceValue
 * @typedef {import('vue').Ref} Ref
 * @typedef {import('vue').ToRefs} ToRefs
 * @typedef {import('vue').MaybeRefOrGetter} MaybeRefOrGetter
 */

/**
 * @param {MaybeRefOrGetter<string|PullAirQualityRequest.AsObject>} query - The name of the device or a query object
 * @param {MaybeRefOrGetter<boolean>=} paused - Whether to pause the data stream
 * @return {ToRefs<ResourceValue<AirQuality.AsObject, PullAirQualityResponse>>}
 */
export function usePullAirQuality(query, paused = false) {
  const resource = reactive(
      /** @type {ResourceValue<AirQuality.AsObject, PullAirQualityResponse>} */
      newResourceValue()
  );
  onScopeDispose(() => closeResource(resource));

  const queryObject = computed(() => toQueryObject(query));

  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      (req) => {
        pullAirQualitySensor(req, resource);
        return () => closeResource(resource);
      }
  );

  return toRefs(resource);
}

/**
 * @typedef {Object} AirQualityMetricDesc
 * @property {string} label
 * @property {string} unit
 * @property {number?} min
 * @property {number?} max
 * @property {Array<{value: number, status: string}>} levels
 */
/**
 * @type {Record<string, AirQualityMetricDesc>}
 */
export const metrics = {
  'score': {
    label: 'IAQ',
    unit: '%',
    min: 0,
    max: 100,
    levels: [
      {value: 0, status: 'error'},
      {value: 10, status: 'warning'},
      {value: 50, status: 'success'}
    ]
  },
  'carbonDioxideLevel': {
    label: 'CO<sub>2</sub>',
    unit: 'ppm',
    min: 0,
    max: 5000,
    levels: [
      {value: 0, status: 'success'},
      {value: 1000, status: 'warning'},
      {value: 2000, status: 'error'}
    ]
  },
  'volatileOrganicCompounds': {
    label: 'VOC',
    unit: 'ppm',
    min: 0,
    max: 1,
    levels: [
      {value: 0, status: 'success'},
      {value: 0.3, status: 'warning'},
      {value: 0.5, status: 'error'}
    ]
  },
  'airPressure': {
    label: 'Air Pressure',
    unit: 'hPa',
    min: 0,
    max: 1100,
    levels: [
      {value: 0, status: 'error'},
      {value: 1000, status: 'success'}
    ]
  },
  'infectionRisk': {
    label: 'Infection Risk',
    unit: '%',
    min: 0,
    max: 100,
    levels: [
      {value: 0, status: 'success'},
      {value: 25, status: 'warning'},
      {value: 50, status: 'error'}
    ]
  },
  'particulateMatter1': {
    label: 'PM1',
    unit: 'µg/m³',
    min: 0,
    max: 50,
    levels: [
      {value: 0, status: 'success'},
      {value: 10, status: 'warning'},
      {value: 20, status: 'error'}
    ]
  },
  'particulateMatter25': {
    label: 'PM2.5',
    unit: 'µg/m³',
    min: 0,
    max: 50,
    levels: [
      {value: 0, status: 'success'},
      {value: 10, status: 'warning'},
      {value: 20, status: 'error'}
    ]
  },
  'particulateMatter10': {
    label: 'PM10',
    unit: 'µg/m³',
    min: 0,
    max: 50,
    levels: [
      {value: 0, status: 'success'},
      {value: 10, status: 'warning'},
      {value: 20, status: 'error'}
    ]
  },
  'airChangePerHour': {
    label: 'Air Exchange Rate',
    unit: '/h',
    min: 0,
    max: 10,
    levels: [
      {value: 0, status: 'error'},
      {value: 5, status: 'success'}
    ]
  },
  'comfort': {
    label: 'Comfort',
    unit: '',
    levels: [
      {value: AirQuality.Comfort.COMFORTABLE, status: 'success'},
      {value: AirQuality.Comfort.UNCOMFORTABLE, status: 'error'}
    ]
  }
}

/**
 * @typedef {Object} AirQualityMetric
 * @property {number} value
 * @property {string} status
 */

/**
 * @typedef {Object} AirQualityScore
 * @property {number} value
 * @property {string} label
 * @property {string} status
 */

/**
 * @param {MaybeRefOrGetter<AirQuality.AsObject>} value
 * @return {{
 *   presentMetrics: import('vue').ComputedRef<Record<keyof metrics, AirQualityMetric>>,
 *   score: import('vue').ComputedRef<AirQualityScore>
 * }}
 */
export function useAirQuality(value) {
  const _v = computed(() => toValue(value));

  const presentMetrics = computed(() => {
    const result = /** @type {Record<keyof metrics, AirQualityMetric>} */ {};
    for (const [k, v] of Object.entries(_v.value ?? {})) {
      const m = metrics[k];
      if (!m || !v) {
        continue; // skip unknown metrics
      }
      let status = '';
      for (const level of m.levels) {
        if (v >= level.value) {
          status = level.status;
          continue;
        }
        break;
      }
      result[k] = {value: v, status};
    }
    return result;
  });

  const score = computed(() => {
    const presentScore = presentMetrics.value['score'];
    const statusToLabel = (status) => {
      switch (status) {
        case 'error': return 'Poor';
        case 'warning': return 'Fair';
        case 'success': return 'Good';
        default: return '';
      }
    }
    if (presentScore) {
      return {
        value: presentScore.value,
        label: statusToLabel(presentScore.status),
        status: presentScore.status
      };
    }

    return null;
  });

  return {
    presentMetrics,
    score
  };
}
