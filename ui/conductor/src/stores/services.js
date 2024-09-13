import {listServices, pullServices} from '@/api/ui/services';
import {serviceName} from '@/util/gateway';
import {Collection} from '@/util/query';
import {defineStore} from 'pinia';
import {reactive, ref} from 'vue';

export const useServicesStore = defineStore('services', () => {
  const servicesCollections =
      reactive(/** @type {Map<string, Collection>} */{});

  /**
   * @typedef {Object} Service
   * @property {Collection} servicesCollection
   */
  /**
   * @param {string} service
   * @param {string} address
   * @param {string} controllerName
   * @return {Service}
   */
  function getService(service, address = '', controllerName = '') {
    if (!servicesCollections.hasOwnProperty(address)) servicesCollections[address] = {};
    const _serviceName = serviceName(controllerName, service);
    if (!servicesCollections[address].hasOwnProperty(_serviceName)) {
      servicesCollections[address][_serviceName] = newServicesCollection(controllerName);
    }
    return {
      servicesCollection: servicesCollections[address][_serviceName]
    };
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
    getService
  };
});
