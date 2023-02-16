import {defineStore} from 'pinia';
import {newActionTracker} from '@/api/resource';
import {reactive, ref} from 'vue';
import {getServiceMetadata, listServices, pullServices} from '@/api/ui/services';
import {Collection} from '@/util/query';

export const useServicesStore = defineStore('services', () => {
  const metadataTrackers = reactive(/** @type {Map<string, ActionTracker<ServiceMetadata.AsObject>>} */{});
  const servicesCollections =
      reactive(/** @type {Map<string, Collection>} */{});


  /**
   * @typedef service
   * @param {ActionTracker<ServiceMetadata.AsObject>} metadataTracker
   * @param {Collection} servicesCollection
   */
  /**
   * @param {string} service
   * @return {service}
   */
  function getService(service) {
    if (!metadataTrackers.hasOwnProperty(service)) metadataTrackers[service] = newActionTracker();
    if (!servicesCollections.hasOwnProperty(service)) servicesCollections[service] = newServicesCollection();
    return {
      metadataTracker: metadataTrackers[service],
      servicesCollection: servicesCollections[service]
    };
  }

  /**
   * @param {string} service
   */
  async function refreshMetadata(service) {
    await getServiceMetadata({name: service}, getService(service).metadataTracker);
  }

  /**
   *
   * @return {Collection}
   */
  function newServicesCollection() {
    const listFn = async (name, tracker, pageToken, recordFn) => {
      const page = await listServices({name, pageToken, pageSize: 100}, tracker);
      for (const service of page.servicesList) {
        service.config = JSON.parse(service.configRaw);
        recordFn(service, service.id);
      }
      return page.nextPageToken;
    };
    const pullFn = (name, resources) => {
      pullServices({name}, resources);
    };
    return new Collection(listFn, pullFn);
  }

  return {
    getService,
    refreshMetadata,
    newServicesCollection
  };
});
