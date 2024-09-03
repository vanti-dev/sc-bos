import {defineAsyncComponent} from 'vue';

export const builtinLayouts = {
  'LayoutMainSide': defineAsyncComponent(() => import('@/layout/LayoutMainSide.vue'))
};
