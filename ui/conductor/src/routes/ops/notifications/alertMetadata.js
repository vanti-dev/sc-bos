import {newResourceValue} from '@/api/resource';
import {pullAlertMetadata} from '@/api/ui/alerts';
import {useHubStore} from '@/stores/hub';
import {useUiConfigStore} from '@/stores/ui-config';
import {convertProtoMap} from '@/util/proto';
import {defineStore} from 'pinia';
import {computed, reactive} from 'vue';

/** @typedef {import('@sc-bos/ui-gen/proto/alerts_pb').AlertMetadata} AlertMetadata */

export const useAlertMetadata = defineStore('alertMetadata', () => {
  const alertMetadata = reactive(
      /** @type {ResourceValue<AlertMetadata.AsObject, AlertMetadata>} */ newResourceValue()
  );
  const uiConfig = useUiConfigStore();
  const hubStore = useHubStore();

  /**
   * @return {Promise}
   */
  function init() {
    // wait for config to load
    return uiConfig.configPromise.then((config) => {
      if (config.gateway) {
        // wait for hub info to load
        hubStore.hubPromise
            .then((hub) => {
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

  const acknowledgedCountMap = computed(() => convertProtoMap(alertMetadata.value?.acknowledgedCountsMap));
  const resolvedCountMap = computed(() => convertProtoMap(alertMetadata.value?.resolvedCountsMap));
  const floorCountsMap = computed(() => convertProtoMap(alertMetadata.value?.floorCountsMap));
  const zoneCountsMap = computed(() => convertProtoMap(alertMetadata.value?.zoneCountsMap));
  const subsystemCountsMap = computed(() => convertProtoMap(alertMetadata.value?.subsystemCountsMap));
  const severityCountsMap = computed(() => convertProtoMap(alertMetadata.value?.severityCountsMap));
  const needsAttentionCountsMap = computed(() => convertProtoMap(alertMetadata.value?.needsAttentionCountsMap));

  const badgeCount = computed(() => needsAttentionCountsMap.value['nack_unresolved']);
  const unacknowledgedAlertCount = computed(() => acknowledgedCountMap.value[false]);

  const alertError = computed(() => alertMetadata.streamError);

  return {
    alertMetadata,

    init,

    // Return 0 when the total count is not known
    totalCount: computed(() => (alertMetadata.value?.totalCount ?? 0)),
    acknowledgedCountMap,
    resolvedCountMap,
    floorCountsMap,
    zoneCountsMap,
    subsystemCountsMap,
    severityCountsMap,
    needsAttentionCountsMap,

    badgeCount,
    unacknowledgedAlertCount,

    alertError
  };
});
