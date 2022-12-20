import {routeTitle} from '@/util/router.js';
import {computed} from 'vue';
import {useRoute} from 'vue-router/composables';

/**
 *
 * @return {*}
 */
export function usePage() {
  const currentRoute = useRoute();
  const pageTitle = computed(() => {
    if (!currentRoute) return undefined;
    return routeTitle(currentRoute);
  });

  const hasSections = computed(() => currentRoute?.matched?.some(r => r.components?.sections));
  const hasNav = computed(() => currentRoute?.matched?.some(r => r.components?.nav));

  return {pageTitle, hasSections, hasNav};
}
