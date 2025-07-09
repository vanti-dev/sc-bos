import {defineAsyncComponent} from 'vue';

export const builtinLayouts = {
  'LayoutMainSide': defineAsyncComponent(() => import('@/dynamic/layout/LayoutMainSide.vue')),
  'LayoutGrid': defineAsyncComponent(() => import('@/dynamic/layout/LayoutGrid.vue')),
  // full pages, not widget containers
  'page/SecurityEvents': defineAsyncComponent(() => import('@/routes/ops/security-events/SecurityEventsTable.vue')),
  'page/Waste': defineAsyncComponent(() => import('@/routes/ops/waste/WasteTable.vue'))
};
