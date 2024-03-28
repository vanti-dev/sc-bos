import {newResourceValue} from '@/api/resource.js';
import {pullAirQualitySensor} from '@/api/sc/traits/air-quality-sensor';
import {camelToSentence} from '@/util/string.js';
import {toQueryObject, watchResource} from '@/util/traits.js';
import {toValue} from '@/util/vue';
import {AirQuality} from '@smart-core-os/sc-api-grpc-web/traits/air_quality_sensor_pb';
import {computed, reactive} from 'vue';

/**
 * @param {MaybeRefOrGetter<string|PullAirQualityRequest.AsObject>} query - The name of the device or a query object
 * @param {MaybeRefOrGetter<boolean>} paused - Whether to pause the data stream
 * @return {{
 *  airQualityValue: ResourceValue<AirQuality.AsObject, PullAirQualityResponse>,
 *  airQualityHasScore: import('vue').ComputedRef<boolean>,
 *  airQualityScore: import('vue').ComputedRef<number>,
 *  airQualityScoreColor: import('vue').ComputedRef<string>,
 *  airQualityInformation: import('vue').ComputedRef<{[key: string]: string}|{}>,
 *  error: import('vue').ComputedRef<ResourceError>,
 *  loading: import('vue').ComputedRef<boolean>
 * }}
 */
export default function(query, paused) {
  const airQualityValue = reactive(
      /** @type {ResourceValue<AirQuality.AsObject, PullAirQualityResponse>} */
      newResourceValue()
  );

  const queryObject = computed(() => toQueryObject(query));

  // Utility function to call the API with the query and the resource
  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      (req) => {
        pullAirQualitySensor(req, airQualityValue);
        return airQualityValue;
      }
  );

  // ---------------- AirQuality Values ---------------- //
  /** @type {import('vue').ComputedRef<boolean>} */
  const airQualityHasScore = computed(() => {
    return airQualityValue.value?.score !== undefined;
  });

  /** @type {import('vue').ComputedRef<number>} */
  const airQualityScore = computed(() => {
    return airQualityValue.value?.score || 0;
  });


  // ---------------- AirQuality Styles ---------------- //
  /** @type {import('vue').ComputedRef<string>} */
  const airQualityScoreColor = computed(() => {
    if (airQualityScore.value < 10) {
      return 'error lighten-1';
    } else if (airQualityScore.value < 50) {
      return 'warning';
    } else if (airQualityScore.value < 75) {
      return 'secondary';
    } else {
      return 'success lighten-1';
    }
  });

  // ---------------- AirQuality Information ---------------- //
  /** @type {import('vue').ComputedRef<{[key: string]: string}|{}>} */
  const airQualityInformation = computed(() => {
    if (airQualityValue.value) {
      const data = {};
      Object.entries(airQualityValue.value).forEach(([key, value]) => {
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
    }
    return {};
  });

  // ---------------- Error ---------------- //
  /** @type {import('vue').ComputedRef<ResourceError>} */
  const error = computed(() => {
    return airQualityValue.streamError;
  });

  // ---------------- Loading ---------------- //
  /** @type {import('vue').ComputedRef<boolean>} */
  const loading = computed(() => {
    return airQualityValue.loading;
  });

  return {
    airQualityValue,

    airQualityHasScore,
    airQualityScore,
    airQualityScoreColor,
    airQualityInformation,

    error,
    loading
  };
}
