import {grpcWebEndpoint} from '@/api/config';
import {closeResource, newActionTracker, newResourceCollection} from '@/api/resource';
import {listHubNodes, pullHubNodes} from '@/api/sc/traits/hub';
import {useAppConfigStore} from '@/stores/app-config';
import {defineStore} from 'pinia';
import {computed, reactive, set, watch} from 'vue';

export const useHubStore = defineStore('hub', () => {
  const appConfig = useAppConfigStore();
  const nodesListCollection = reactive(newResourceCollection());

  watch(() => appConfig.config, async config => {
    closeResource(nodesListCollection);

    if (config?.hub) {
      pullHubNodes(nodesListCollection);
      try {
        const nodes = await listHubNodes(newActionTracker());
        for (const node of nodes.nodesList) {
          set(nodesListCollection.value, node.name, node);
        }
      } catch (e) {
        console.warn('Error fetching first page', e);
      }
    }
  }, {immediate: true});

  const nodesList = computed(() => {
    const nodes = {};
    Object.values(nodesListCollection?.value || {}).forEach((node, name) => {
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
    nodesList
  };
});
