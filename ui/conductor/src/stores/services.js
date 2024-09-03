import {newActionTracker} from '@/api/resource';
import {getServiceMetadata, listServices, pullServices} from '@/api/ui/services';
import {serviceName} from '@/util/gateway';
import {Collection} from '@/util/query';
import {defineStore} from 'pinia';
import {reactive, ref} from 'vue';

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
   * @param {string} controllerName
   * @return {Service}
   */
  function getService(service, address = '', controllerName = '') {
    if (!metadataTrackers.hasOwnProperty(address)) metadataTrackers[address] = {};
    if (!servicesCollections.hasOwnProperty(address)) servicesCollections[address] = {};
    const _serviceName = serviceName(controllerName, service);
    if (!metadataTrackers[address].hasOwnProperty(_serviceName)) {
      metadataTrackers[address][_serviceName] = newActionTracker();
    }
    if (!servicesCollections[address].hasOwnProperty(_serviceName)) {
      servicesCollections[address][_serviceName] = newServicesCollection(controllerName);
    }
    return {
      metadataTracker: metadataTrackers[address][_serviceName],
      servicesCollection: servicesCollections[address][_serviceName]
    };
  }

  /**
   * @param {string} service
   * @param {string} address
   * @param {string} controllerName
   * @return {Promise<ServiceMetadata.AsObject>}
   */
  function refreshMetadata(service, address = '', controllerName = '') {
    return getServiceMetadata(
        {name: serviceName(controllerName, service)},
        getService(service, address, controllerName).metadataTracker
    );
  }

  /**
   *
   * @param {string} controllerName
   * @return {Collection}
   */
  function newServicesCollection(controllerName = '') {
    const listFn = async (name, tracker, pageToken, recordFn) => {
      if (name) {
        const page = await listServices({
          name: serviceName(controllerName, name),
          pageToken, pageSize: 100
        }, tracker);
        for (const service of page.servicesList) {
          service.config = JSON.parse(service.configRaw);
          recordFn(service, service.id);
        }
        return page.nextPageToken;
      }
    };
    const pullFn = (name, resources) => {
      if (name) {
        pullServices({name: serviceName(controllerName, name)}, resources);
      }
    };
    return new Collection(listFn, pullFn);
  }

  // The SC node we're getting services from, if absent get services from the node we're communicating directly with.
  const node = ref(null);

  return {
    node,
    getService,
    refreshMetadata,
    newServicesCollection
  };
});
