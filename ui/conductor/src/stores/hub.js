import {grpcWebEndpoint} from '@/api/config';
import {closeResource, newActionTracker, newResourceCollection} from '@/api/resource';
import {getEnrollment} from '@/api/sc/traits/enrollment';
import {listHubNodes, pullHubNodes} from '@/api/sc/traits/hub';
import {listServices, ServiceNames} from '@/api/ui/services';
import {useControllerStore} from '@/stores/controller';
import {useUiConfigStore} from '@/stores/ui-config';
import {isGatewayId} from '@/util/gateway';
import {defineStore} from 'pinia';
import {computed, reactive, ref, watch} from 'vue';

export const useHubStore = defineStore('hub', () => {
  const uiConfig = useUiConfigStore();
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

  watch(() => uiConfig.config, async config => {
    closeResource(nodesListCollection);

    if (config?.hub) {
      pullHubNodes(nodesListCollection);
      await nodesListCollectionInit();
    }
  }, {immediate: true});

  const nodesListCollectionInit = async () => {
    try {
      // if local gateway hub mode is enabled, the hub node will be the same as the gateway node
      // get systems config, so we can check if the gateway is in local mode
      const systems = await listServices({name: ServiceNames.Systems}, newActionTracker());
      let gatewayHubLocalMode = false;

      // search through systems to find the gateway
      for (const system of systems.servicesList) {
        if (isGatewayId(system.id)) {
          const cfg = JSON.parse(system.configRaw);
          // check hub mode
          if (cfg.hubMode && cfg.hubMode === 'local') {
            gatewayHubLocalMode = true;
            break;
          }
        }
      }

      if (gatewayHubLocalMode) {
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
      _hubReject(e);
    }
  };


  const listHubNodesAction = async (actionTracker) => {
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
  };
  /**
   * @typedef {Object} Node
   * @property {string} name - the Smart Core name of the node
   * @property {string} address - the address of the node
   * @property {string} description - a human-readable description of the node
   * @property {string} commsAddress - the address to use to communicate with the node (based on gateway settings)
   * @property {string} commsName - the name to use to communicate with the node (based on gateway settings)
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
   * If we're communicating with named devices via a gateway, this will return prepend the controller name to the
   * resource, otherwise it will return the resource name as-is.
   *
   * @param {string} name
   * @param {string} address
   * @return {string}
   */
  async function proxiedName(name, address) {
    // check if running in gateway mode, and that the node address is not the same as our endpoint address
    if (uiConfig.config?.gateway && (await grpcWebEndpoint()) !== address) {
      return name;
    }
    return '';
  }

  /**
   * If we're communicating with named devices via a gateway, this will return the gateway address, otherwise it will
   * return the address as-is.
   *
   * @param {string} address
   * @return {string}
   */
  async function proxiedAddress(address) {
    if (uiConfig.config?.gateway) {
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
