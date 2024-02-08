import {defineStore} from 'pinia';
import {useUiConfigStore} from '@/stores/ui-config';
import {computed} from 'vue';

export const useWidgetsStore = defineStore('widgets', () => {
  const {config, defaultConfig} = useUiConfigStore();

  // Returns the active building overview widgets
  // We merge the default and the imported config values
  // If the config value is not set, the default value will be used,
  // otherwise any existing config value will override the default value
  const activeOverviewWidgets = computed(() => {
    return {
      ...defaultConfig.config.ops.overview.widgets,
      ...config?.ops?.overview?.widgets
    };
  });

  return {
    activeOverviewWidgets
  };
});
