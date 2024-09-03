import {defineAsyncComponent} from 'vue';

export const builtinWidgets = {
  'environmental/EnvironmentalCard': defineAsyncComponent(
      () => import('@/widgets/environmental/EnvironmentalCard.vue')),
  'graphic/LayeredGraphic': defineAsyncComponent(() => import('@/widgets/graphic/LayeredGraphic.vue')),
  'occupancy/OccupancyCard': defineAsyncComponent(() => import('@/widgets/occupancy/OccupancyCard.vue')),
  'power-history/PowerHistoryCard': defineAsyncComponent(() => import('@/widgets/power-history/PowerHistoryCard.vue')),
  'occupancy/PresenceCard': defineAsyncComponent(() => import('@/widgets/occupancy/PresenceCard.vue')),
  'notifications/ZoneNotifications': defineAsyncComponent(() => import('@/widgets/notifications/ZoneNotifications.vue'))
};
