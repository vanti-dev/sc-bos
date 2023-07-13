import {closeResource, newResourceValue} from '@/api/resource';
import {pullAlertMetadata} from '@/api/ui/alerts';
import {useErrorStore} from '@/components/ui-error/error';
import {useAppConfigStore} from '@/stores/app-config';
import {useHubStore} from '@/stores/hub';
import {convertProtoMap} from '@/util/proto';
import {defineStore} from 'pinia';
import {computed, onMounted, onUnmounted, reactive} from 'vue';

/** @typedef {import('@sc-bos/ui-gen/proto/alerts_pb').AlertMetadata} AlertMetadata */

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
        hubStore.hubPromise
            .then(hub => {
              pullAlertMetadata({name: hub.name, updatesOnly: false}, alertMetadata);
            })
            .catch(() => {
              pullAlertMetadata({name: '', updatesOnly: false}, alertMetadata);
            });
      } else {
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
  const resolvedCountMap = computed(() => convertProtoMap(alertMetadata.value?.resolvedCountsMap));
  const floorCountsMap = computed(() => convertProtoMap(alertMetadata.value?.floorCountsMap));
  const zoneCountsMap = computed(() => convertProtoMap(alertMetadata.value?.zoneCountsMap));
  const severityCountsMap = computed(() => convertProtoMap(alertMetadata.value?.severityCountsMap));
  const needsAttentionCountsMap = computed(() => convertProtoMap(alertMetadata.value?.needsAttentionCountsMap));

  const badgeCount = computed(() => needsAttentionCountsMap.value['nack_unresolved']);
  const unacknowledgedAlertCount = computed(() => acknowledgedCountMap.value[false]);

  return {
    alertMetadata,

    init,

    totalCount: computed(() => alertMetadata.value?.totalCount),
    acknowledgedCountMap,
    resolvedCountMap,
    floorCountsMap,
    zoneCountsMap,
    severityCountsMap,
    needsAttentionCountsMap,

    badgeCount,
    unacknowledgedAlertCount
  };
});
