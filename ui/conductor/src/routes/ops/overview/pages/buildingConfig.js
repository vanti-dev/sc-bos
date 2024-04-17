import {useUiConfigStore} from '@/stores/ui-config.js';
import {useWidgetsStore} from '@/stores/widgets.js';
import {computed} from 'vue';

/**
 * @return {{
 *   energyZone: ComputedRef<string>,
 *   occupancyZone: ComputedRef<string>,
 *   environmentalZone: ComputedRef<string>,
 *   showEnergy: ComputedRef<boolean>,
 *   showEnvironment: ComputedRef<boolean>,
 *   externalZone: ComputedRef<string>,
 *   buildingZone: ComputedRef<string>,
 *   supplyZone: ComputedRef<string>,
 *   showEnergyChart: ComputedRef<boolean>,
 *   showEnergyIntensity: ComputedRef<boolean>,
 *   showOccupancy: ComputedRef<boolean>
 * }}
 */
export default function useBuildingConfig() {
  const uiConfig = useUiConfigStore();
  const {activeOverviewWidgets} = useWidgetsStore();

  const buildingZoneSource = computed(() => uiConfig.config?.ops?.buildingZone ?? '');
  const supplyZoneSource = computed(() => uiConfig.config?.ops?.supplyZone ?? '');

  const buildingZone = computed(() => buildingZoneSource.value);
  const supplyZone = computed(() => supplyZoneSource.value ? supplyZoneSource.value + '/supply' : '');
  const energyZone = buildingZone;
  const showEnergy = computed(() => {
    if (typeof activeOverviewWidgets.showEnergyConsumption === 'boolean') {
      return activeOverviewWidgets.showEnergyConsumption;
    } else {
      return activeOverviewWidgets.showEnergyConsumption?.showChart ||
          activeOverviewWidgets.showEnergyConsumption?.showIntensity;
    }
  });
  const showEnergyChart = computed(() => {
    if (typeof activeOverviewWidgets.showEnergyConsumption === 'boolean') {
      return activeOverviewWidgets.showEnergyConsumption;
    } else {
      return activeOverviewWidgets.showEnergyConsumption?.showChart;
    }
  });
  const showEnergyIntensity = computed(() => {
    if (typeof activeOverviewWidgets.showEnergyConsumption === 'boolean') {
      return activeOverviewWidgets.showEnergyConsumption;
    } else {
      return activeOverviewWidgets.showEnergyConsumption?.showIntensity;
    }
  });

  const showOccupancy = computed(() => activeOverviewWidgets?.showOccupancy);
  const occupancyZone = buildingZone;

  const showEnvironment = computed(() => activeOverviewWidgets?.showEnvironment);
  const environmentalZone = buildingZone;
  const externalZone = computed(() => environmentalZone.value + '/outside');

  return {
    buildingZone,
    supplyZone,
    energyZone,
    showEnergy,
    showEnergyChart,
    showEnergyIntensity,
    showOccupancy,
    occupancyZone,
    showEnvironment,
    environmentalZone,
    externalZone
  };
}
