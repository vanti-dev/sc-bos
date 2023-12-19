import {closeResource, newActionTracker, newResourceValue} from '@/api/resource';
import {pullBrightness, updateBrightness} from '@/api/sc/traits/light';
import {useErrorStore} from '@/components/ui-error/error';
import {onMounted, onUnmounted, reactive, watch} from 'vue';

/**
 *
 * @param {Object} props
 * @param {string} props.name
 * @param {boolean} [props.paused]
 * @return {{
 *  lightValue: import('vue').UnwrapNestedRefs<
 *    ResourceValue<Brightness.AsObject, Brightness>
 *   >,
 *  updateValue: import('vue').UnwrapNestedRefs<
 *   ActionTracker<Brightness.AsObject>
 *  >,
 *  doUpdateBrightness: (req: number|Brightness.AsObject|UpdateBrightnessRequest.AsObject) => void
 * }}
 */
export default function(props) {
  const errorStore = useErrorStore();

  const lightValue = reactive(
      /** @type {ResourceValue<Brightness.AsObject, Brightness>} */
      newResourceValue());

  const updateValue = reactive(
      /** @type {ActionTracker<Brightness.AsObject>}  */
      newActionTracker());


  //
  //
  // Methods
  /**
   * @param {number|Brightness.AsObject|UpdateBrightnessRequest.AsObject} req
   */
  function doUpdateBrightness(req) {
    if (typeof req === 'number') {
      req = {levelPercent: Math.min(100, Math.round(req))};
    }
    if (!req.hasOwnProperty('brightness')) {
      req = {brightness: req};
    }
    req.name = props.name;
    updateBrightness(req, updateValue);
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
          closeResource(lightValue);
        }

        if (!newPaused && (oldPaused || newName !== oldName)) {
          closeResource(lightValue);
          pullBrightness({name: newName}, lightValue);
        }
      },
      {immediate: true, deep: true, flush: 'sync'}
  );

  //
  //
  // UI error handling
  let unwatchUpdateError;
  onMounted(() => {
    unwatchUpdateError = errorStore.registerTracker(updateValue);
  });
  onUnmounted(() => {
    closeResource(lightValue);
    if (unwatchUpdateError) unwatchUpdateError();
  });

  return {
    lightValue,
    updateValue,
    doUpdateBrightness
  };
}