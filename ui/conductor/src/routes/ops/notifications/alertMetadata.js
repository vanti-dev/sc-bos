import {closeResource, newResourceValue} from '@/api/resource';
import {pullAlertMetadata} from '@/api/ui/alerts';
import {useControllerStore} from '@/stores/controller';
import {defineStore} from 'pinia';
import {computed, reactive, watch} from 'vue';

export const useAlertMetadata = defineStore('alertMetadata', () => {
  const controller = useControllerStore();
  const alertMetadata = reactive(/** @type {ResourceValue<AlertMetadata.AsObject, AlertMetadata>} */newResourceValue());
  watch(() => controller.controllerName, async name => {
    closeResource(alertMetadata);
    pullAlertMetadata({name, updatesOnly: false}, alertMetadata);
  }, {immediate: true});

  /**
   * Converts a proto map, which is an array of [k,v] into a js object.
   *
   * @param {Array<[K,V]>} arr
   * @return {Object<K,V>}
   * @template K,V
   */
  function convertProtoMap(arr) {
    if (!arr) return {};
    const dst = {};
    for (const [k, v] of arr || []) {
      dst[k] = v;
    }
    return dst;
  }

  const acknowledgedCountMap = computed(() => convertProtoMap(alertMetadata.value?.acknowledgedCountsMap));
  const floorCountsMap = computed(() => convertProtoMap(alertMetadata.value?.floorCountsMap));
  const zoneCountsMap = computed(() => convertProtoMap(alertMetadata.value?.zoneCountsMap));
  const severityCountsMap = computed(() => convertProtoMap(alertMetadata.value?.severityCountsMap));

  const unacknowledgedAlertCount = computed(() => acknowledgedCountMap.value[false]);

  return {
    alertMetadata,

    acknowledgedCountMap,
    floorCountsMap,
    zoneCountsMap,
    severityCountsMap,

    unacknowledgedAlertCount
  };
});
