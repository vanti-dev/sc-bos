import {defineStore} from 'pinia';
import {newActionTracker} from '@/api/resource';
import {computed, reactive} from 'vue';
import {getServiceMetadata, ServiceNames} from '@/api/ui/services';

export const useAutomationsStore = defineStore('automations', () => {
  const metadataTracker = reactive(/** @type {ActionTracker<ServiceMetadata.AsObject>} */newActionTracker());

  // filter out automations that have no instances, and map to {type, number} obj
  const automationList = computed(() => {
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

  return {
    automationList,
    refreshMetadata
  };
});
