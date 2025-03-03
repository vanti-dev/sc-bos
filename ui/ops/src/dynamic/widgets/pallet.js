import {defineAsyncComponent} from 'vue';

export const builtinWidgets = {
  'environmental/EnvironmentalCard': defineAsyncComponent(
      () => import('@/dynamic/widgets/environmental/EnvironmentalCard.vue')),
  'graphic/LayeredGraphic': defineAsyncComponent(() => import('@/dynamic/widgets/graphic/LayeredGraphic.vue')),
  'occupancy/OccupancyCard': defineAsyncComponent(() => import('@/dynamic/widgets/occupancy/OccupancyCard.vue')),
  'power-history/PowerHistoryCard': defineAsyncComponent(() => import('@/dynamic/widgets/power-history/PowerHistoryCard.vue')),
  'occupancy/PresenceCard': defineAsyncComponent(() => import('@/dynamic/widgets/occupancy/PresenceCard.vue')),
  'notifications/ZoneNotifications': defineAsyncComponent(() => import('@/dynamic/widgets/notifications/ZoneNotifications.vue')),
  // from elsewhere in our codebase
  'lighting/LightIcon': defineAsyncComponent(() => import('@/traits/light/LightIcon.vue')),
  'meter/ConsumptionCard': defineAsyncComponent(() => import('@/traits/meter/ConsumptionCard.vue')),
  'environmental/AirTemperatureChip': defineAsyncComponent(() => import('@/traits/airTemperature/AirTemperatureChip.vue')),
  'devices/DeviceTable': defineAsyncComponent(() => import('@/routes/devices/components/DeviceTable.vue')),
  // containers
  'container/FlexRow': defineAsyncComponent(() => import('@/dynamic/widgets/container/FlexRow.vue')),
};
