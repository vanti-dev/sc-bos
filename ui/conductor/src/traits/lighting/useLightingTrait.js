import {newActionTracker, newResourceValue} from '@/api/resource';
import {pullBrightness, updateBrightness} from '@/api/sc/traits/light';
import {toQueryObject, watchResource} from '@/util/traits.js';
import {toValue} from '@/util/vue';
import {computed, reactive} from 'vue';

/**
 *
 * @template T
 * @param {MaybeRefOrGetter<T>} query - The name of the device or a query object
 * @param {MaybeRefOrGetter<boolean>} paused - Whether to pause the data stream
 * @param {{
 *  pullBrightness?: boolean,
 *  describeBrightness?: boolean
 * }} [options] - Options to control which API calls to make
 * @return {{
 *  brightnessValue: ResourceValue<Brightness.AsObject, Brightness>,
 *  brightnessUpdate: ActionTracker<Brightness.AsObject>,
 *  brightnessLevelNumber: import('vue').ComputedRef<number>,
 *  brightnessLevelString: import('vue').ComputedRef<string>,
 *  lightIcon: import('vue').ComputedRef<string>,
 *  updateBrightness: (req: number|Brightness.AsObject|UpdateBrightnessRequest.AsObject) => void,
 *  error: import('vue').ComputedRef<ResourceError|ActionError>,
 *  loading: import('vue').ComputedRef<boolean>
 * }}
 */
export default function(query, paused, options) {
  const brightnessValue = reactive(
      /** @type {ResourceValue<Brightness.AsObject, Brightness>} */
      newResourceValue());

  const brightnessUpdate = reactive(
      /** @type {ActionTracker<Brightness.AsObject>}  */
      newActionTracker());

  /**
   * Create a list of API calls to make based on the options passed in
   *
   * @type {Array<(req: any) => any>} apiCalls
   */
  const apiCalls = [];

  // For light state and updates
  if (options?.pullBrightness !== false) {
    apiCalls.push((req) => {
      pullBrightness(req, brightnessValue);
      return brightnessValue;
    });
  }

  const queryObject = toQueryObject(query); // Make sure the query is an object

  // Utility function to call the API with the query and the resource
  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      ...apiCalls
  );
  //
  //
  // ---------------- Light and Brightness Display ---------------- //
  const brightnessLevelNumber = computed(() => brightnessValue.value?.levelPercent);

  const brightnessLevelString = computed(() => {
    if (brightnessLevelNumber.value === 0) {
      return 'Off';
    } else if (brightnessLevelNumber.value === 100) {
      return 'Max';
    } else if (brightnessLevelNumber.value > 0 && brightnessLevelNumber.value < 100) {
      return `${brightnessLevelNumber.value.toFixed(0)}%`;
    }

    return '';
  });

  const lightIcon = computed(() => {
    if (brightnessLevelNumber.value === 0) {
      return 'mdi-lightbulb-outline';
    } else if (brightnessLevelNumber.value > 0) {
      return 'mdi-lightbulb-on';
    } else {
      return '';
    }
  });


  // ---------------- Light and Brightness Control ---------------- //
  /**
   * @param {number|Brightness.AsObject|UpdateBrightnessRequest.AsObject} request
   * @return {void}
   */
  const doBrightnessUpdate = (request) => {
    if (typeof request === 'number') {
      request = {
        levelPercent: Math.min(100, Math.round(request))
      };
    }
    if (!request.hasOwnProperty('brightness')) {
      request = {
        brightness: request
      };
    }
    request.name = toValue(queryObject).name;

    updateBrightness(request, brightnessUpdate);
  };


  // ----------- Errors ----------- //
  const error = computed(() => {
    return brightnessValue.streamError || brightnessUpdate.error;
  });


  // ----------- Loading ----------- //
  const loading = computed(() => {
    return brightnessValue.loading || brightnessUpdate.loading;
  });

  // ---------------- Return ---------------- //
  return {
    // Resources
    brightnessValue,
    brightnessUpdate,

    // Computed
    brightnessLevelNumber,
    brightnessLevelString,
    lightIcon,

    // Actions
    updateBrightness: doBrightnessUpdate,

    // States
    error,
    loading
  };
}
