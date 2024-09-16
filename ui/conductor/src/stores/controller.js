import {usePullMetadata} from '@/traits/metadata/metadata.js';
import {acceptHMRUpdate, defineStore} from 'pinia';
import {computed, watch} from 'vue';

/**
 * A store that reports on the server (controller) the ui is directly connected to.
 */
export const useControllerStore = defineStore('controller', () => {
  // fetch the metadata for the server we're connected to
  const {value, streamError, loading} = usePullMetadata({});
  const controllerName = computed(() => value.value?.name);
  const controllerNameError = computed(() => {
    const error = streamError.value?.error;
    if (error) error.from = 'useControllerStore';
    return error;
  });
  const hasLoaded = computed(() => Boolean(!loading.value && (streamError.value || value.value)));

  let notifyLoaded;
  const waitForLoad = new Promise((resolve) => notifyLoaded = resolve);
  watch(hasLoaded, (loaded) => {
    if (loaded) notifyLoaded();
  }, {immediate: true});

  return {
    controllerName,
    hasLoaded,
    controllerNameError,
    waitForLoad
  };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useControllerStore, import.meta.hot));
}
