import {defineStore} from 'pinia';
import {newActionTracker, newResourceCollection} from '@/api/resource';
import {computed, reactive} from 'vue';
import {getServiceMetadata, listServices, ServiceNames} from '@/api/ui/services';
import {Collection} from '@/util/query';

export const useAutomationsStore = defineStore('automations', () => {
  const metadataTracker = reactive(/** @type {ActionTracker<ServiceMetadata.AsObject>} */newActionTracker());
  const automationsCollection = reactive(
      /** @type {ResourceCollection<Service.AsObject, Service>} */newResourceCollection());

  // filter out automations that have no instances, and map to {type, number} obj
  const automationTypeList = computed(() => {
    if (!metadataTracker.response) return [];
    const list = [];
    metadataTracker.response.typeCountsMap.forEach(([type, number]) => {
      if (number > 0) {
        list.push({type, number});
      }
    });
    return list;
  });

  /**
   *
   */
  async function refreshMetadata() {
    await getServiceMetadata({name: ServiceNames.Automations}, metadataTracker);
  }

  /**
   *
   * @return {Collection}
   */
  function newAutomationsCollection() {
    const listFn = async (name, tracker, pageToken, recordFn) => {
      const page = await listServices({name, pageToken, pageSize: 100}, tracker);
      for (const service of page.servicesList) {
        recordFn(service, service.id);
      }
      return page.nextPageToken;
    };
    return new Collection(listFn);
  }

  return {
    automationTypeList,
    automationsCollection,
    refreshMetadata,
    newAutomationsCollection
  };
});
