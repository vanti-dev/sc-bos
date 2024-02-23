import {newActionTracker, newResourceValue} from '@/api/resource';
import {describeBrightness, pullBrightness, updateBrightness} from '@/api/sc/traits/light';
import {toQueryObject, watchResource} from '@/util/traits.js';
import {toValue} from '@/util/vue';
import {computed, reactive} from 'vue';

/**
 * @typedef {Object} LevelPercent
 * @property {number} levelPercent
 */

/**
 * @typedef {Object} Preset
 * @property {LightPreset.AsObject} preset
 */

/**
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
 *  brightnessSupport: ActionTracker<BrightnessSupport.AsObject>,
 *  brightnessLevelNumber: import('vue').ComputedRef<number>,
 *  brightnessLevelString: import('vue').ComputedRef<string>,
 *  brightnessPresets: import('vue').ComputedRef<Array<Preset.preset>>,
 *  lightIcon: import('vue').ComputedRef<string>,
 *  toLevelPercentObject: (percent: LevelPercent|number) => LevelPercent,
 *  toPresetObject: (preset: Brightness.AsObject.LightPreset.AsObject|Preset.preset) => Preset,
 *  updateBrightness: (
 *    req: {LevelPercent}|Preset|Brightness.AsObject|UpdateBrightnessRequest.AsObject
 *  ) => void,
 *  error: import('vue').ComputedRef<ResourceError|ActionError>,
 *  loading: import('vue').ComputedRef<boolean>
 * }}
 */
export default function(query, paused, options) {
  const brightnessValue = reactive(/** @type {ResourceValue<Brightness.AsObject, Brightness>} */
      newResourceValue());

  const brightnessUpdate = reactive(/** @type {ActionTracker<Brightness.AsObject>}  */
      newActionTracker());

  const brightnessSupport = reactive(/** @type {ActionTracker<BrightnessSupport.AsObject>}  */
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
  if (options?.describeBrightness !== false) {
    apiCalls.push((req) => {
      describeBrightness(req, brightnessSupport);
      return brightnessSupport;
    });
  }

  const queryObject = toQueryObject(query); // Make sure the query is an object

  // Utility function to call the API with the query and the resource
  watchResource(() => toValue(queryObject), () => toValue(paused), ...apiCalls);
  //
  //
  // ---------------- Light and Brightness Display ---------------- //

  /** @type {import('vue').ComputedRef<number>} */
  const brightnessLevelNumber = computed(() => brightnessValue.value?.levelPercent || 0);

  /** @type {import('vue').ComputedRef<string>} */
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

  /** @type {import('vue').ComputedRef<Array<Preset.preset>>} */
  const brightnessPresets = computed(() => {
    if (brightnessSupport.response?.presetsList) {
      return brightnessSupport.response.presetsList;
    } else {
      return [];
    }
  });

  /** @type {import('vue').ComputedRef<string>} */
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
   * Convert a number to a levelPercent object if it is not already
   *
   * @param {number | LevelPercent} percent
   * @return {LevelPercent}
   */
  const toLevelPercentObject = (percent) => {
    if (typeof percent === 'number') {
      return {
        levelPercent: Math.min(100, Math.round(percent))
      };
    } else {
      return percent;
    }
  };

  /**
   * Convert a preset to a preset object if it is not already
   *
   * @param {Brightness.AsObject.LightPreset.AsObject | LightPreset.AsObject} preset
   * @return {Preset}
   */
  const toPresetObject = (preset) => {
    if (!preset.hasOwnProperty('preset')) {
      return {preset};
    } else {
      return preset;
    }
  };

  /**
   * @param {{LevelPercent}|Preset|Brightness.AsObject|UpdateBrightnessRequest.AsObject} request
   * @return {void}
   */
  const doBrightnessUpdate = (request) => {
    let modifiedRequest = request;

    // Convert the request to an UpdateBrightnessRequest
    if (!request.hasOwnProperty('brightness')) {
      modifiedRequest = {
        brightness: request
      };
    }

    // Set the name of the device
    modifiedRequest.name = toValue(queryObject).name;

    // Call the API
    updateBrightness(modifiedRequest, brightnessUpdate);
  };


  // ----------- Errors ----------- //

  /** @type {import('vue').ComputedRef<ResourceError|ActionError>} */
  const error = computed(() => {
    return brightnessValue.streamError || brightnessUpdate.error || brightnessSupport.error;
  });


  // ----------- Loading ----------- //

  /** @type {import('vue').ComputedRef<boolean>} */
  const loading = computed(() => {
    return brightnessValue.loading || brightnessUpdate.loading || brightnessSupport.loading;
  });

  // ---------------- Return ---------------- //
  return {
    // Resources
    brightnessValue,
    brightnessUpdate,
    brightnessSupport,

    // Computed
    brightnessLevelNumber,
    brightnessLevelString,
    brightnessPresets,
    lightIcon,

    // Utilities
    toLevelPercentObject,
    toPresetObject,

    // Actions
    updateBrightness: doBrightnessUpdate,

    // States
    error,
    loading
  };
}
