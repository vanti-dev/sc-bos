import {newActionTracker, newResourceValue} from '@/api/resource';
import {describeBrightness, pullBrightness, updateBrightness} from '@/api/sc/traits/light';
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
 *  light: ResourceValue<Brightness.AsObject, Brightness>,
 *  lightUpdate: ActionTracker<Brightness.AsObject>,
 *  brightnessSupport: ActionTracker<BrightnessSupport.AsObject>,
 *  brightnessLevelNumber: import('vue').ComputedRef<number>,
 *  brightnessLevelString: import('vue').ComputedRef<string>,
 *  lightIcon: import('vue').ComputedRef<string>,
 *  lightPresets: import('vue').ComputedRef<{name: string, title: string}[]>,
 *  lightError: import('vue').ComputedRef<ResourceError|ActionError>,
 *  presetsError: import('vue').ComputedRef<ResourceError>,
 *  doUpdateBrightness: (req: number|Brightness.AsObject|UpdateBrightnessRequest.AsObject) => void,
 *  lightLoading: import('vue').ComputedRef<boolean>,
 *  presetsLoading: import('vue').ComputedRef<boolean>
 * }}
 */
export default function(query, paused, options) {
  const light = reactive(
      /** @type {ResourceValue<Brightness.AsObject, Brightness>} */
      newResourceValue());

  const lightUpdate = reactive(
      /** @type {ActionTracker<Brightness.AsObject>}  */
      newActionTracker());

  const brightnessSupport = reactive(
      /** @type {ActionTracker<BrightnessSupport.AsObject>}  */
      newActionTracker());

  // Make sure the query is an object
  const queryObject = toQueryObject(query);

  /**
   * Create a list of API calls to make based on the options passed in
   *
   * @type {import('vue').ComputedRef<Array<(req: any) => any>>} apiCalls
   */
  const apiCalls = computed(() => [
    // For light state and updates
    ...(options?.pullBrightness !== false ? [(req) => {
      pullBrightness(req, light);
      return light;
    }] : []),

    // For brightness presets and support
    ...(options?.describeBrightness !== false ? [(req) => {
      describeBrightness(req, brightnessSupport);
      return brightnessSupport;
    }] : [])
  ]);

  // Utility function to call the API with the query and the resource
  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      ...apiCalls.value
  );
  //
  //
  // ---------------- Light Display ---------------- //
  const brightnessLevelNumber = computed(() => light.value?.levelPercent);

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

  const lightPresets = computed(() => {
    if (brightnessSupport.response?.value?.presets) {
      return brightnessSupport.response.value.presets;
    } else {
      return [];
    }
  });


  // ---------------- Light Control ---------------- //
  /**
   * @param {number|Brightness.AsObject|UpdateBrightnessRequest.AsObject} request
   * @return {void}
   */
  const doUpdateBrightness = (request) => {
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

    updateBrightness(request, lightUpdate);
  };


  // ----------- Errors ----------- //
  const lightError = computed(() => {
    return light.streamError || lightUpdate.error;
  });

  const presetsError = computed(() => {
    return brightnessSupport.error;
  });


  // ----------- Loading ----------- //
  const lightLoading = computed(() => {
    return light.loading || lightUpdate.loading;
  });

  const presetsLoading = computed(() => {
    return brightnessSupport.loading;
  });

  // ---------------- Return ---------------- //
  return {
    light,
    lightUpdate,
    brightnessSupport,

    brightnessLevelNumber,
    brightnessLevelString,
    lightIcon,
    lightPresets,

    doUpdateBrightness,
    lightError,
    presetsError,

    lightLoading,
    presetsLoading
  };
}
