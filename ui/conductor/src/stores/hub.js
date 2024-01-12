import {grpcWebEndpoint} from '@/api/config';
import {closeResource, newActionTracker, newResourceCollection} from '@/api/resource';
import {getEnrollment} from '@/api/sc/traits/enrollment';
import {listHubNodes, pullHubNodes} from '@/api/sc/traits/hub';
import {listServices, ServiceNames} from '@/api/ui/services';
import {useAppConfigStore} from '@/stores/app-config';
import {useControllerStore} from '@/stores/controller';
import {defineStore} from 'pinia';
import {computed, reactive, ref, watch} from 'vue';

export const useHubStore = defineStore('hub', () => {
  const appConfig = useAppConfigStore();
  const controller = useControllerStore();
  const nodesListCollection = reactive(newResourceCollection());
  const hubNode = ref();
  const listedHubNodes = ref([]);
  let _hubResolve;
  let _hubReject;
  const hubPromise = new Promise((resolve, reject) => {
    _hubResolve = resolve;
    _hubReject = reject;
  });

  watch(() => appConfig.config, async config => {
    closeResource(nodesListCollection);

    if (config?.hub) {
      pullHubNodes(nodesListCollection);
      await nodesListCollectionInit();
    }
  }, {immediate: true});

  const nodesListCollectionInit = async () => {
    try {
      // if local proxy hub mode is enabled, the hub node will be the same as the proxy node
      // get systems config, so we can check if the proxy is in local mode
      const systems = await listServices({name: ServiceNames.Systems}, newActionTracker());
      let proxyHubLocalMode = false;

      // search through systems to find the proxy
      for (const system of systems.servicesList) {
        if (system.id === 'proxy') {
          const cfg = JSON.parse(system.configRaw);
          // check hub mode
          if (cfg.hubMode && cfg.hubMode === 'local') {
            proxyHubLocalMode = true;
            break;
          }
        }
      }

      if (proxyHubLocalMode) {
        hubNode.value = {
          name: controller.controllerName,
          address: ''
        };
      } else {
        const hub = await getEnrollment(newActionTracker());
        hubNode.value = {
          name: hub.managerName,
          address: hub.managerAddress
        };
      }

      // add the hub node to the list
      nodesListCollection.value = {
        ...nodesListCollection.value,
        [hubNode.value.name]: hubNode.value
      };


      await listHubNodesAction(newActionTracker());
      _hubResolve(hubNode.value);
    } catch (e) {
      console.warn('Error fetching first page', e);
      _hubReject(e);
    }
  };


  const listHubNodesAction = async (actionTracker) => {
    try {
      const nodes = await listHubNodes(actionTracker);

      // reset the existing list
      listedHubNodes.value = [];

      // add the new nodes to the list
      for (const node of nodes.nodesList) {
        // collect the names of the nodes
        listedHubNodes.value.push(node.name);

        // updating the reactive object while keeping the reactivity
        nodesListCollection.value = {
          ...nodesListCollection.value,
          [node.name]: node
        };
      }
    } catch (error) {
      console.error('Error in listHubNodesAction:', error);
      throw error;
    }
  };


  /**
   * @typedef {Object} Node
   * @property {string} name - the Smart Core name of the node
   * @property {string} address - the address of the node
   * @property {string} description - a human-readable description of the node
   * @property {string} commsAddress - the address to use to communicate with the node (based on proxy settings)
   * @property {string} commsName - the name to use to communicate with the node (based on proxy settings)
   */

  const nodesList = computed(() => {
    /** @type {Record<string, Node>} */
    const nodes = {};
    Object.values(nodesListCollection?.value || {}).forEach((node) => {
      nodes[node.name] = {
        ...node,
        commsAddress: proxiedAddress(node.address),
        commsName: proxiedName(node.name, node.address)
      };
    });
    return nodes;
  });


  /**
   * If we're communicating with named devices via a proxy, this will return prepend the controller name to the
   * resource, otherwise it will return the resource name as-is.
   *
   * @param {string} name
   * @param {string} address
   * @return {string}
   */
  async function proxiedName(name, address) {
    // check if running in proxy mode, and that the node address is not the same as our endpoint address
    if (appConfig.config?.proxy && (await grpcWebEndpoint()) !== address) {
      return name;
    }
    return '';
  }

  /**
   * If we're communicating with named devices via a proxy, this will return the proxy address, otherwise it will return
   * the address as-is.
   *
   * @param {string} address
   * @return {string}
   */
  async function proxiedAddress(address) {
    if (appConfig.config?.proxy) {
      return await grpcWebEndpoint();
    }
    return address;
  }

  return {
    listedHubNodes,
    nodesList,
    hubNode,
    hubPromise,
    nodesListCollection,
    listHubNodesAction
  };
});
