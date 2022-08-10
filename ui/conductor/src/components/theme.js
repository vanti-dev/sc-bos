import {computed} from 'vue';
import vuetify from '../plugins/vuetify.js';
import {useRoute} from '../util/router.js';

export function useTheme() {
  const currentRoute = useRoute();
  const themeName = computed(() => currentRoute?.name || '');
  const logoColor = computed(() => {
    return currentRoute?.meta?.['logoColor'] || vuetify.framework.theme.themes.dark[themeName.value] || undefined;
  })

  return {logoColor, themeName};
}
