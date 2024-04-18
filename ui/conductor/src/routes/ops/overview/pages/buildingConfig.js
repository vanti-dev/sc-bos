import {builtinLayouts} from '@/layout/pallet.js';
import {useUiConfigStore} from '@/stores/ui-config.js';
import {isNullOrUndef} from '@/util/types.js';
import {toValue} from '@/util/vue.js';
import {builtinWidgets} from '@/widgets/pallet.js';
import {computed, ref} from 'vue';

/**
 * @return {{
 *   layout: import('vue').Component,
 *   title: Ref<string>,
 *   main: ComputedRef<{component: Component, props: Object}[]>,
 *   after: ComputedRef<{component: Component, props: Object}[]>
 * }}
 */
export default function useBuildingConfig() {
  const uiConfig = useUiConfigStore();

  const buildingZone = computed(
      () => /** @type {string|undefined} */ uiConfig.getOrDefault('ops.buildingZone'));
  const supplyZone = computed(() => {
    const sz = /** @type {string|undefined} */ uiConfig.getOrDefault('ops.supplyZone');
    // for legacy reasons we append /supply to the configured value
    return isNullOrUndef(sz) ? sz : sz + '/supply';
  });

  const demandSource = buildingZone;
  const generatedSource = supplyZone;
  const occupancySource = buildingZone;
  const envInternalSource = buildingZone;
  const envExternalSource = computed(() => envInternalSource.value + '/outside');

  const showPowerHistory = computed(() => {
    const v = uiConfig.getOrDefault('ops.overview.widgets.showEnergyConsumption');
    if (typeof v === 'boolean') {
      return /** @type {boolean} */ v;
    } else {
      return uiConfig.getOrDefault('ops.overview.widgets.showEnergyConsumption.showChart') ||
          uiConfig.getOrDefault('ops.overview.widgets.showEnergyConsumption.showIntensity');
    }
  });
  const powerHistoryConfig = computed(() => {
    if (!showPowerHistory.value) return false;
    return {
      demand: demandSource.value,
      generated: generatedSource.value,
      occupancy: occupancySource.value,
      hideChart: /** @type {boolean} */
          !uiConfig.getOrDefault('ops.overview.widgets.showEnergyConsumption.showChart', true),
      hideTotal: /** @type {boolean} */
          !uiConfig.getOrDefault('ops.overview.widgets.showEnergyConsumption.showIntensity', true)
    };
  });

  const showOccupancy = computed(
      () => uiConfig.getOrDefault('ops.overview.widgets.showOccupancy', true));
  const occupancyHistoryConfig = computed(() => {
    if (!showOccupancy.value) return false;
    return {
      source: occupancySource.value
    };
  });

  const showEnvironment = computed(
      () => uiConfig.getOrDefault('ops.overview.widgets.showEnvironment', true));
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
    layout: builtinLayouts['LayoutMainSide'],
    title: ref('Building Status Overview'),
    main: computed(() => {
      const res = [];
      addIfPresent(res, powerHistoryConfig, builtinWidgets['power-history/PowerHistoryCard']);
      addIfPresent(res, occupancyHistoryConfig, builtinWidgets['occupancy/OccupancyCard']);
      return res;
    }),
    after: computed(() => {
      const res = [];
      addIfPresent(res, environmentalConfig, builtinWidgets['environmental/EnvironmentalCard']);
      return res;
    })
  };
}
