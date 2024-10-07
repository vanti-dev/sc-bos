import {defineAsyncComponent} from 'vue';

export const builtinLayouts = {
  'LayoutMainSide': defineAsyncComponent(() => import('@/dynamic/layout/LayoutMainSide.vue'))
};
