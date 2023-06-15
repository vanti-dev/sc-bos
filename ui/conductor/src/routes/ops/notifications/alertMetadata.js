import {closeResource, newResourceValue} from '@/api/resource';
import {pullAlertMetadata} from '@/api/ui/alerts';
import {useErrorStore} from '@/components/ui-error/error';
import {useAppConfigStore} from '@/stores/app-config';
import {useHubStore} from '@/stores/hub';
import {convertProtoMap} from '@/util/proto';
import {defineStore} from 'pinia';
import {computed, onMounted, onUnmounted, reactive} from 'vue';

export const useAlertMetadata = defineStore('alertMetadata', () => {
  const alertMetadata = reactive(/** @type {ResourceValue<AlertMetadata.AsObject, AlertMetadata>} */newResourceValue());
  const appConfig = useAppConfigStore();
  const hubStore = useHubStore();


  /**
   * @return {Promise}
   */
  function init() {
    // wait for config to load
    return appConfig.configPromise.then(config => {
      if (config.proxy) {
        // wait for hub info to load
        hubStore.hubPromise.then(hub => {
          console.debug('Fetching alert metadata for', hub.name);
          pullAlertMetadata({name: hub.name, updatesOnly: false}, alertMetadata);
        });
      } else {
        console.debug('Fetching alert metadata for current node');
        pullAlertMetadata({name: '', updatesOnly: false}, alertMetadata);
      }
    });
  }

  // Ui Error Handling
  const errorStore = useErrorStore();
  let unwatchErrors;
  onMounted(() => {
    unwatchErrors = errorStore.registerValue(alertMetadata);
  });
  onUnmounted(() => {
    closeResource(alertMetadata);
    if (unwatchErrors) unwatchErrors();
  });

  const acknowledgedCountMap = computed(() => convertProtoMap(alertMetadata.value?.acknowledgedCountsMap));
  const floorCountsMap = computed(() => convertProtoMap(alertMetadata.value?.floorCountsMap));
  const zoneCountsMap = computed(() => convertProtoMap(alertMetadata.value?.zoneCountsMap));
  const severityCountsMap = computed(() => convertProtoMap(alertMetadata.value?.severityCountsMap));

  const unacknowledgedAlertCount = computed(() => acknowledgedCountMap.value[false]);

  return {
    alertMetadata,

    init,

    acknowledgedCountMap,
    floorCountsMap,
    zoneCountsMap,
    severityCountsMap,

    unacknowledgedAlertCount
  };
});
