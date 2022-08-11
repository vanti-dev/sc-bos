import {computed} from 'vue';
import vuetify from '../plugins/vuetify.js';
import {routeTitle, useRoute} from '../util/router.js';

export function usePage() {
  const currentRoute = /** @type {import('vue-router').Route} */ useRoute();
  const themeColor = computed(() => {
    if (!currentRoute) return undefined;
    for (let i = currentRoute.matched.length - 1; i >= 0; i--) {
      const r = currentRoute.matched[i];
      const color = r.meta?.['logoColor'] || vuetify.framework.theme.currentTheme[r.name];
      if (color) return color;
    }
  });
  const pageTitle = computed(() => {
    if (!currentRoute) return undefined;
    return routeTitle(currentRoute);
  });

  const hasSections = computed(() => currentRoute?.matched?.some(r => r.components?.sections));

  return {themeColor, pageTitle, hasSections};
}
