import {closeResource, newResourceValue} from '@/api/resource.js';
import {pullAirQualitySensor} from '@/api/sc/traits/air-quality-sensor';
import {camelToSentence} from '@/util/string.js';
import {toQueryObject, watchResource} from '@/util/traits.js';
import {toValue} from '@/util/vue';
import {AirQuality} from '@smart-core-os/sc-api-grpc-web/traits/air_quality_sensor_pb';
import {computed, onUnmounted, reactive, toRefs} from 'vue';

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
 * @typedef {import('@/util/vue').MaybeRefOrGetter} MaybeRefOrGetter
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
  onUnmounted(() => closeResource(resource));

  const queryObject = computed(() => toQueryObject(query));

  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      (req) => {
        pullAirQualitySensor(req, resource);
        return resource;
      }
  );

  return toRefs(resource);
}

/**
 * @param {MaybeRefOrGetter<AirQuality.AsObject>} value
 * @return {{
 *   hasScore: ComputedRef<boolean>,
 *   score: ComputedRef<number>,
 *   scoreColor: ComputedRef<string>,
 *   tableData: ComputedRef<{[key: string]: string}>
 * }}
 */
export function useAirQuality(value) {
  const _v = computed(() => toValue(value));

  const hasScore = computed(() => _v.value?.score !== undefined);
  const score = computed(() => _v.value?.score ?? 0);
  const scoreColor = computed(() => {
    if (score.value < 10) {
      return 'error lighten-1';
    } else if (score.value < 50) {
      return 'warning';
    } else if (score.value < 75) {
      return 'secondary';
    } else {
      return 'success lighten-1';
    }
  });

  const tableData = computed(() => {
    const data = {};
    Object.entries(_v.value ?? {}).forEach(([key, value]) => {
      if (value !== undefined) {
        switch (key) {
          case 'carbonDioxideLevel':
            if (value > 0) {
              data['CO2'] = Math.round(value) + ' ppm';
            }
            break;
          case 'volatileOrganicCompounds':
            if (value > 0) {
              data['VOC'] = value.toFixed(3) + ' ppm';
            }
            break;
          case 'airPressure':
            if (value > 0) {
              data['Air Pressure'] = Math.round(value) + ' hPa';
            }
            break;
          case 'infectionRisk':
            if (value > 0) {
              data['Infection Risk'] = Math.round(value) + '%';
            }
            break;
          case 'particulateMatter1':
            if (value > 0) {
              data['PM1'] = value.toFixed(1) + ' µg/m³';
            }
            break;
          case 'particulateMatter25':
            if (value > 0) {
              data['PM2.5'] = value.toFixed(1) + ' µg/m³';
            }
            break;
          case 'particulateMatter10':
            if (value > 0) {
              data['PM10'] = value.toFixed(1) + ' µg/m³';
            }
            break;
          case 'airChangePerHour':
            if (value > 0) {
              data['Air Exchange Rate'] = value.toFixed(1) + ' /h';
            }
            break;
          case 'comfort':
            switch (value) {
              case AirQuality.Comfort.COMFORTABLE:
                data['Comfort'] = 'Comfortable';
                break;
              case AirQuality.Comfort.UNCOMFORTABLE:
                data['Comfort'] = 'Uncomfortable';
                break;
              default:
                // do nothing
            }
            break;
          case 'score':
            if (value > 0) {
              data['Air Quality Score'] = Math.round(value) + '%';
            }
            break;
          default: {
            data[camelToSentence(key)] = value;
          }
        }
      }
    });
    return data;
  });

  return {
    hasScore,
    score,
    scoreColor,
    tableData
  };
}
