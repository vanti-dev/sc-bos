import {newActionTracker, newResourceValue} from '@/api/resource';
import {
  airTemperatureModeToString,
  pullAirTemperature,
  temperatureToString,
  updateAirTemperature
} from '@/api/sc/traits/air-temperature';
import {toQueryObject, watchResource} from '@/util/traits';
import {toValue} from '@/util/vue';
import {AirTemperature} from '@smart-core-os/sc-api-grpc-web/traits/air_temperature_pb';
import {computed, reactive, ref} from 'vue';

/**
 * @template T
 * @param {MaybeRefOrGetter<T>} query - The name of the device or a query object
 * @param {MaybeRefOrGetter<boolean>} paused - Whether to pause the data stream
 * @return {{
 *   airTemperatureValue: ResourceValue<AirTemperature.AsObject, PullAirTemperatureResponse>,
 *   airTemperatureUpdate: ActionTracker<AirTemperature.AsObject>,
 *   hasSetPoint: import('vue').ComputedRef<boolean>,
 *   hasTemperature: import('vue').ComputedRef<boolean>,
 *   hasHumidity: import('vue').ComputedRef<boolean>,
 *   ambientTemperature: import('vue').ComputedRef<number>,
 *   airTemperatureSetPointNumber: import('vue').ComputedRef<number>,
 *   airTemperatureSetPointString: import('vue').ComputedRef<string>,
 *   ambientHumidity: import('vue').ComputedRef<number>,
 *   airTemperatureProgress: import('vue').ComputedRef<number>,
 *   airTemperatureInformation: import('vue').ComputedRef<{[key: string]: string|number}|{}>,
 *   toTemperatureSetPointObject: (
 *     valueCelsius: number|Partial<AirTemperature.AsObject>|Partial<UpdateAirTemperatureRequest.AsObject>
 *   ) => Partial<UpdateAirTemperatureRequest.AsObject>,
 *   updateAirTemperature: (
 *     req: Partial<AirTemperature.AsObject>|Partial<UpdateAirTemperatureRequest.AsObject>
 *   ) => void,
 *   error: import('vue').ComputedRef<ResourceError|ActionError>,
 *   loading: import('vue').ComputedRef<boolean>
 * }}
 */
export default function(query, paused) {
  const airTemperatureValue = reactive(
      /** @type {ResourceValue<AirTemperature.AsObject, PullAirTemperatureResponse>} */
      newResourceValue()
  );

  const airTemperatureUpdate = reactive(
      /** @type {ActionTracker<AirTemperature.AsObject>}  */
      newActionTracker()
  );

  const queryObject = computed(() => toQueryObject(query));

  // Utility function to call the API with the query and the resource
  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      (req) => {
        pullAirTemperature(req, airTemperatureValue);
        return airTemperatureValue;
      }
  );

  // --------------------- Air Temperature Values --------------------- //
  /** @type {import('vue').Ref<{low: number, high: number}>} defaultTemperatureRange */
  const defaultTemperatureRange = ref({
    low: 18.0,
    high: 24.0
  });

  /** @type {import('vue').ComputedRef<boolean>} */
  const hasSetPoint = computed(() => {
    return airTemperatureValue.value?.temperatureSetPoint?.valueCelsius !== undefined;
  });

  /** @type {import('vue').ComputedRef<boolean>} */
  const hasTemperature = computed(() => {
    return airTemperatureValue.value?.ambientTemperature?.valueCelsius !== undefined;
  });

  /** @type {import('vue').ComputedRef<boolean>} */
  const hasHumidity = computed(() => {
    return airTemperatureValue.value?.ambientHumidity !== undefined;
  });

  /** @type {import('vue').ComputedRef<number>} */
  const ambientTemperature = computed(() => airTemperatureValue?.value?.ambientTemperature?.valueCelsius);

  /** @type {import('vue').ComputedRef<number>} */
  const airTemperatureSetPointNumber = computed(() => {
    return airTemperatureValue.value?.temperatureSetPoint?.valueCelsius;
  });

  const UNIT = '°C'; // Temperature unit

  /** @type {import('vue').ComputedRef<string>} */
  const airTemperatureSetPointString = computed(() => {
    const formattedTemperature = airTemperatureSetPointNumber.value?.toFixed(1);

    if (hasSetPoint.value) {
      return `${formattedTemperature}`;
    } else {
      return formattedTemperature + UNIT;
    }
  });

  /** @type {import('vue').ComputedRef<number>} */
  const ambientHumidity = computed(() => airTemperatureValue?.value?.ambientHumidity);

  // --------------------- Air Temperature Information --------------------- //

  /** @type {import('vue').ComputedRef<number>} airTemperatureProgress */
  const airTemperatureProgress = computed(() => {
    let val = airTemperatureValue.value?.ambientTemperature?.valueCelsius ?? 0;
    if (val > 0) {
      val -= defaultTemperatureRange.value.low;
      val = val / (defaultTemperatureRange.value.high - defaultTemperatureRange.value.low);
    }
    return val * 100;
  });

  /** @type {import('vue').ComputedRef<{[key: string]: string|number}|{}>} airTemperatureInformation */
  const airTemperatureInformation = computed(() => {
    if (airTemperatureValue.value) {
      const data = {};
      Object.entries(airTemperatureValue.value).forEach(([key, value]) => {
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

  // --------------------- Air Temperature Control --------------------- //
  /**
   * Convert the request to UpdateAirTemperatureRequest partially with the temperatureSetPoint set if it is not already
   *
   * @param {number|Partial<AirTemperature.AsObject>|Partial<UpdateAirTemperatureRequest.AsObject>} valueCelsius
   * @return {Partial<UpdateAirTemperatureRequest.AsObject>}
   */
  const toTemperatureSetPointObject = (valueCelsius) => {
    if (typeof valueCelsius === 'number') {
      return {temperatureSetPoint: {valueCelsius}};
    } else {
      return valueCelsius;
    }
  };

  /**
   * Convert the request to UpdateAirTemperatureRequest partially with the updateMask set if it is not already
   *
   * @param {Partial<AirTemperature.AsObject>|Partial<UpdateAirTemperatureRequest.AsObject>} request
   * @return {Partial<UpdateAirTemperatureRequest.AsObject>}
   */
  const toUpdateMaskObject = (request) => {
    if (!request.hasOwnProperty('updateMask')) {
      return {
        ...request,
        updateMask: {pathsList: ['temperature_set_point']}
      };
    } else {
      return request;
    }
  };

  /**
   * Convert the request to an UpdateAirTemperatureRequest if it is not already
   *
   * @param {Partial<AirTemperature.AsObject>|Partial<UpdateAirTemperatureRequest.AsObject>} state
   * @return {Partial<UpdateAirTemperatureRequest.AsObject>}
   */
  const toStateObject = (state) => {
    if (!state.hasOwnProperty('state')) {
      return toUpdateMaskObject({state});
    } else {
      return toUpdateMaskObject(state);
    }
  };

  /**
   * @param {Partial<AirTemperature.AsObject>|Partial<UpdateAirTemperatureRequest.AsObject>} request
   * @return {void}
   */
  const doAirTemperatureUpdate = (request) => {
    // Convert the request to UpdateAirTemperatureRequest if it is not already
    const modifiedRequest = toStateObject(request);

    // Make sure the name is set
    modifiedRequest.name = toValue(queryObject).name;

    // Call the API with the modified request
    updateAirTemperature(modifiedRequest, airTemperatureUpdate);
  };


  const error = computed(() => {
    return airTemperatureValue.streamError || airTemperatureUpdate.error;
  });

  const loading = computed(() => {
    return airTemperatureValue.loading || airTemperatureUpdate.loading;
  });

  return {
    airTemperatureValue,
    airTemperatureUpdate,
    hasSetPoint,
    hasTemperature,
    hasHumidity,
    ambientTemperature,
    airTemperatureSetPointNumber,
    airTemperatureSetPointString,
    ambientHumidity,
    airTemperatureProgress,
    airTemperatureInformation,
    toTemperatureSetPointObject,
    updateAirTemperature: doAirTemperatureUpdate,
    error,
    loading
  };
}
