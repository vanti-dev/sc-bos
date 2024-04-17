import {useUiConfigStore} from '@/stores/ui-config.js';
import {isNullOrUndef} from '@/util/types.js';
import {toValue} from '@/util/vue.js';
import {builtinWidgets} from '@/widgets/pallet.js';
import {computed, ref} from 'vue';

/**
 * @return {{
 *   title: Ref<string>,
 *   main: ComputedRef<{component: Component, props: Object}[]>,
 *   after: ComputedRef<{component: Component, props: Object}[]>
 * }}
 */
export default function useBuildingConfig() {
  const uiConfig = useUiConfigStore();

  /**
   * Gets the value of path from either uiConfig config or defaultConfig, depending on presence.
   *
   * @template T
   * @param {string} path
   * @param {T?} def
   * @return {T}
   */
  const getOrDefault = (path, def) => {
    const parts = path.split('.');
    let a = uiConfig.config;
    let b = uiConfig.defaultConfig?.config;
    for (let i = 0; i < parts.length; i++) {
      a = a?.[parts[i]];
      b = b?.[parts[i]];
    }
    return a ?? b ?? toValue(def);
  };

  const buildingZone = computed(
      () => /** @type {string|undefined} */ getOrDefault('ops.buildingZone'));
  const supplyZone = computed(() => {
    const sz = /** @type {string|undefined} */ getOrDefault('ops.supplyZone');
    // for legacy reasons we append /supply to the configured value
    return isNullOrUndef(sz) ? sz : sz + '/supply';
  });

  const demandSource = buildingZone;
  const generatedSource = supplyZone;
  const occupancySource = buildingZone;
  const envInternalSource = buildingZone;
  const envExternalSource = computed(() => envInternalSource.value + '/outside');

  const showPowerHistory = computed(() => {
    const v = getOrDefault('ops.overview.widgets.showEnergyConsumption');
    if (typeof v === 'boolean') {
      return /** @type {boolean} */ v;
    } else {
      return getOrDefault('ops.overview.widgets.showEnergyConsumption.showChart') ||
          getOrDefault('ops.overview.widgets.showEnergyConsumption.showIntensity');
    }
  });
  const powerHistoryConfig = computed(() => {
    if (!showPowerHistory.value) return false;
    return {
      demand: demandSource.value,
      generated: generatedSource.value,
      occupancy: occupancySource.value,
      hideChart: /** @type {boolean} */ !getOrDefault('ops.overview.widgets.showEnergyConsumption.showChart', true),
      hideTotal: /** @type {boolean} */ !getOrDefault('ops.overview.widgets.showEnergyConsumption.showIntensity', true)
    };
  });

  const showOccupancy = computed(() => getOrDefault('ops.overview.widgets.showOccupancy', true));
  const occupancyHistoryConfig = computed(() => {
    if (!showOccupancy.value) return false;
    return {
      source: occupancySource.value
    };
  });

  const showEnvironment = computed(() => getOrDefault('ops.overview.widgets.showEnvironment', true));
  const environmentalConfig = computed(() => {
    if (!showEnvironment.value) return false;
    return {
      internal: envInternalSource.value,
      external: envExternalSource.value
    };
  });

  const addIfPresent = (arr, props, component) => {
    props = toValue(props);
    if (props) {
      arr.push({props, component});
    }
  };
  return {
    title: ref('Building Status Overview'),
    main: computed(() => {
      const res = [];
      addIfPresent(res, powerHistoryConfig, builtinWidgets['power-history']);
      addIfPresent(res, occupancyHistoryConfig, builtinWidgets['occupancy-history']);
      return res;
    }),
    after: computed(() => {
      const res = [];
      addIfPresent(res, environmentalConfig, builtinWidgets['environmental']);
      return res;
    })
  };
}
