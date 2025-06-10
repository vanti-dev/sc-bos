import {computed} from 'vue';
import {useRoute} from 'vue-router';

/**
 * A composable to handle signage mode in the application.
 * Signage mode is full screen, zero interaction mode; so no scrolling, clicking, etc.
 *
 * @return {{
 *   enabled: import('vue').ComputedRef<boolean>,
 *   styles: import('vue').ComputedRef<Object|Array|string>
 * }}
 */
export default function useSignage() {
  const route = useRoute();

  const enabled = computed(() => {
    return route.query?.mode === 'signage';
  });

  return {
    enabled,
    styles: computed(() => {
      if (!enabled.value) return {};
      return {
        position: 'fixed',
        top: '0',
        left: '0',
        width: '100vw',
        height: '100vh',
        overflow: 'hidden',
        'z-index': '9999',
        background: 'rgb(var(--v-theme-neutral-darken-1))',
      }
    })
  }
}