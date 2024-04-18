export const builtinWidgets = {
  'environmental/EnvironmentalCard': () => import('@/widgets/environmental/EnvironmentalCard.vue'),
  'occupancy/OccupancyCard': () => import('@/widgets/occupancy/OccupancyCard.vue'),
  'power-history/PowerHistoryCard': () => import('@/widgets/power-history/PowerHistoryCard.vue'),
  'occupancy/PresenceCard': () => import('@/widgets/occupancy/PresenceCard.vue'),
  'notifications/ZoneNotifications': () => import('@/widgets/notifications/ZoneNotifications.vue')
};
