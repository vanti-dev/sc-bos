import {closeResource, newActionTracker, newResourceValue} from '@/api/resource';
import {pullAirTemperature, updateAirTemperature} from '@/api/sc/traits/air-temperature';
import {useErrorStore} from '@/components/ui-error/error';
import {onMounted, onUnmounted, reactive, watch} from 'vue';

/**
 *
 * @param {Object} props
 * @param {string} props.name
 * @param {boolean} [props.paused]
 * @return {{
 *  airTemperatureResource: ResourceValue<AirTemperature.AsObject, PullAirTemperatureResponse>,
 *  updateTracker: ActionTracker<AirTemperature.AsObject>,
 *  doUpdateAirTemperature: (
 *    function(number|Partial<AirTemperature.AsObject>|Partial<UpdateAirTemperatureRequest.AsObject>)
 *  )
 * }}
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
  watch(
      [() => props.paused, () => props.name],
      ([newPaused, newName], [oldPaused, oldName]) => {
        // only for LightSensor
        if (newPaused === oldPaused && newName === oldName) return;

        if (newPaused) {
          closeResource(airTemperatureResource);
        }

        if (!newPaused && (oldPaused || newName !== oldName)) {
          closeResource(airTemperatureResource);
          pullAirTemperature({name: newName}, airTemperatureResource);
        }
      },
      {immediate: true, deep: true, flush: 'sync'}
  );

  //
  //
  // UI error handling
  const errorHandlers = [];
  onMounted(() => {
    errorHandlers.push(
        errorStore.registerTracker(updateTracker)
    );
  });
  onUnmounted(() => {
    closeResource(airTemperatureResource);
    errorHandlers.forEach(unwatch => unwatch());
  });

  return {
    airTemperatureResource,
    updateTracker,
    doUpdateAirTemperature
  };
}
