import {defineAsyncComponent} from 'vue';

export const builtinWidgets = {
  'container/FlexRow': defineAsyncComponent(() => import('@/dynamic/widgets/container/FlexRow.vue')),
  'energy/EnergyHistoryCard': defineAsyncComponent(() => import('@/dynamic/widgets/energy/EnergyHistoryCard.vue')),
  'environmental/EnvironmentalCard': defineAsyncComponent(() => import('@/dynamic/widgets/environmental/EnvironmentalCard.vue')),
  'general/CohortStatus': defineAsyncComponent(() => import('@/dynamic/widgets/general/CohortStatus.vue')),
  'general/DateAndTime': defineAsyncComponent(() => import('@/dynamic/widgets/general/DateAndTime.vue')),
  'general/PlaceholderCard': defineAsyncComponent(() => import('@/dynamic/widgets/general/PlaceholderCard.vue')),
  'graphic/LayeredGraphic': defineAsyncComponent(() => import('@/dynamic/widgets/graphic/LayeredGraphic.vue')),
  'notifications/ZoneNotifications': defineAsyncComponent(() => import('@/dynamic/widgets/notifications/ZoneNotifications.vue')),
  'occupancy/OccupancyCard': defineAsyncComponent(() => import('@/dynamic/widgets/occupancy/OccupancyCard.vue')),
  'occupancy/PresenceCard': defineAsyncComponent(() => import('@/dynamic/widgets/occupancy/PresenceCard.vue')),
  'power-history/PowerHistoryCard': defineAsyncComponent(() => import('@/dynamic/widgets/power-history/PowerHistoryCard.vue')),
  // from elsewhere in our codebase
  'devices/DeviceTable': defineAsyncComponent(() => import('@/routes/devices/components/DeviceTable.vue')),
  'environmental/AirTemperatureChip': defineAsyncComponent(() => import('@/traits/airTemperature/AirTemperatureChip.vue')),
  'lighting/LightIcon': defineAsyncComponent(() => import('@/traits/light/LightIcon.vue')),
  'meter/ConsumptionCard': defineAsyncComponent(() => import('@/traits/meter/ConsumptionCard.vue')),
};
