export const builtinWidgets = {
  'environmental': () => import('@/widgets/environmental/EnvironmentalCard.vue'),
  'occupancy-history': () => import('@/widgets/occupancy/OccupancyCard.vue'),
  'power-history': () => import('@/widgets/power-history/PowerHistoryCard.vue')
};
