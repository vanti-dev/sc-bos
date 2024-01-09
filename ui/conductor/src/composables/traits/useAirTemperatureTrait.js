import {closeResource, newActionTracker, newResourceValue} from '@/api/resource';
import {pullAirTemperature, updateAirTemperature} from '@/api/sc/traits/air-temperature';
import {useErrorStore} from '@/components/ui-error/error';
import {watchResource} from '@/util/traits';
import {computed, reactive} from 'vue';

/**
 * @typedef {Object} AirTemperatureTrait
 * @property {ResourceValue<AirTemperature.AsObject, PullAirTemperatureResponse>} airTemperatureResource
 * @property {ActionTracker<AirTemperature.AsObject>} updateTracker
 * @property {
 *  (number|Partial<AirTemperature.AsObject>|Partial<UpdateAirTemperatureRequest.AsObject>)
 * } doUpdateAirTemperature
 * @property {import('vue').ComputedRef<number>} temperatureValue
 * @property {import('vue').ComputedRef<number>} humidityValue
 * @property {function} collectErrors
 * @property {function} clearResourceError
 * @param {Object} props
 * @param {string} props.name
 * @param {boolean} [props.paused]
 * @return {AirTemperatureTrait}
 */
export default function(props) {
  const errorStore = useErrorStore();

  const airTemperatureResource = reactive(
      /** @type {ResourceValue<AirTemperature.AsObject, PullAirTemperatureResponse>} */
      newResourceValue());

  const updateTracker = reactive(
      /** @type {ActionTracker<AirTemperature.AsObject>}  */
      newActionTracker());


  //
  //
  // Methods
  /**
   * @param {number|Partial<AirTemperature.AsObject>|Partial<UpdateAirTemperatureRequest.AsObject>} req
   */
  function doUpdateAirTemperature(req) {
    if (typeof req === 'number') {
      req = {
        state: {temperatureSetPoint: {valueCelsius: /** @type {number} */ req}},
        updateMask: {pathsList: ['temperature_set_point']}
      };
    }
    if (!req.hasOwnProperty('state')) {
      req = {state: /** @type {AirTemperature.AsObject} */ req};
    }
    req.name = props.name;
    updateAirTemperature(req, updateTracker);
  }

  //
  //
  // Watch
  // Depending on paused state/device name, we close/open data stream(s)
  watchResource(
      [() => props.paused, () => props.name],
      airTemperatureResource,
      (params, resource) => {
        pullAirTemperature(params, resource);
      }
  );

  //
  //
  // Return the temperature of the single device specified
  const temperatureValue = computed(() =>
    airTemperatureResource.value?.ambientTemperature?.valueCelsius ?? 0
  );

  // Return the humidity of the single device specified
  const humidityValue = computed(() =>
    airTemperatureResource.value?.ambientHumidity ?? 0
  );

  //
  //
  // UI error handling
  const errorHandlers = [];

  const collectErrors = () => {
    errorHandlers.push(
        errorStore.registerTracker(updateTracker)
    );
  };

  const clearResourceError = () => {
    closeResource(airTemperatureResource);
    errorHandlers.forEach(unwatch => unwatch());
  };

  return {
    airTemperatureResource,
    updateTracker,
    doUpdateAirTemperature,
    temperatureValue,
    humidityValue,
    collectErrors,
    clearResourceError
  };
}
