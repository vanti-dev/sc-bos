import {newActionTracker} from '@/api/resource';
import {getServiceMetadata, listServices, pullServices} from '@/api/ui/services';
import {serviceName} from '@/util/proxy';
import {Collection} from '@/util/query';
import {defineStore} from 'pinia';
import {reactive} from 'vue';

export const useServicesStore = defineStore('services', () => {
  const metadataTrackers =
      reactive(/** @type {Map<string, ActionTracker<ServiceMetadata.AsObject>>} */{});
  const servicesCollections =
      reactive(/** @type {Map<string, Collection>} */{});


  /**
   * @typedef {Object} Service
   * @property {ActionTracker<ServiceMetadata.AsObject>} metadataTracker
   * @property {Collection} servicesCollection
   */
  /**
   * @param {string} service
   * @param {string} address
   * @param {string} name
   * @return {Service}
   */
  function getService(service, address= '', name = '') {
    if (!metadataTrackers.hasOwnProperty(address)) metadataTrackers[address] = {};
    if (!servicesCollections.hasOwnProperty(address)) servicesCollections[address] = {};
    if (!metadataTrackers[address].hasOwnProperty(service)) metadataTrackers[address][service] = newActionTracker();
    if (!servicesCollections[address].hasOwnProperty(service)) {
      servicesCollections[address][service] = newServicesCollection(name);
    }
    return {
      metadataTracker: metadataTrackers[address][service],
      servicesCollection: servicesCollections[address][service]
    };
  }

  /**
   * @param {string} service
   * @param {string} address
   */
  async function refreshMetadata(service, address='') {
    await getServiceMetadata({name: service}, getService(service, address).metadataTracker);
  }

  /**
   *
   * @param {string} controllerName
   * @return {Collection}
   */
  function newServicesCollection(controllerName = '') {
    const listFn = async (name, tracker, pageToken, recordFn) => {
      const page = await listServices({name: serviceName(controllerName, name),
        pageToken, pageSize: 100}, tracker);
      for (const service of page.servicesList) {
        service.config = JSON.parse(service.configRaw);
        recordFn(service, service.id);
      }
      return page.nextPageToken;
    };
    const pullFn = (name, resources) => {
      pullServices({name: serviceName(controllerName, name)}, resources);
    };
    return new Collection(listFn, pullFn);
  }

  return {
    getService,
    refreshMetadata,
    newServicesCollection
  };
});
