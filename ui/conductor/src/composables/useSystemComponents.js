import {computed, onUnmounted, reactive, set, watch} from 'vue';
import {closeResource, newActionTracker} from '@/api/resource';
import {enrollHubNode, forgetHubNode, testHubNode} from '@/api/sc/traits/hub';
import {ServiceNames} from '@/api/ui/services';
import {useEnrollmentStore} from '@/stores/enrollmentStore';
import {useHubStore} from '@/stores/hub';
import {useServicesStore} from '@/stores/services';

/**
 * @return {{
 * nodesList: Object,
 * isProxy: function(string): boolean,
 * isHub: function(string): boolean
 * }}
 */
export default function() {
  const hubStore = useHubStore();
  const servicesStore = useServicesStore();
  const {enrollmentValue} = useEnrollmentStore();


  // --------------------------- //
  // Test Hub Nodes
  const hubNodeValue = reactive(newActionTracker());


  watch(enrollmentValue, (newValue) => {
    if (newValue?.response?.targetAddress) {
      const request = {
        address: newValue.response.targetAddress
      };

      testHubNode(request, hubNodeValue);
    }
  }, {immediate: true, deep: true});


  // --------------------------- //
  // Manage Hub Nodes
  const enrollHubNodeValue = reactive(newActionTracker());

  /**
   * @param {string} address
   * @return {Promise<void>|undefined}
   */
  function enrollHubNodeAction(address) {
    if (!address) return;

    const request = {
      node: {
        address
      },
      publicCertsList: []
    };

    return enrollHubNode(request, enrollHubNodeValue);
  }

  // --------------------------- //
  // List and Track Hub Nodes
  const nodeDetails = reactive({});

  let unwatchTrackers = [];

  const nodesList = computed(() => {
    return Object.values(hubStore.nodesList).map(node => {
      Promise.all([node.commsAddress, node.commsName])
          .then(([address, name]) => {
            console.debug('node', node);
            set(nodeDetails, node.name, {
              automations: servicesStore.getService(ServiceNames.Automations, address, name),
              drivers: servicesStore.getService(ServiceNames.Drivers, address, name),
              systems: servicesStore.getService(ServiceNames.Systems, address, name)
            });
            unwatchTrackers.push(nodeDetails[node.name].automations.metadataTracker);
            unwatchTrackers.push(nodeDetails[node.name].drivers.metadataTracker);
            unwatchTrackers.push(nodeDetails[node.name].systems.metadataTracker);
            return Promise.all([
              servicesStore.refreshMetadata(ServiceNames.Automations, address, name),
              servicesStore.refreshMetadata(ServiceNames.Drivers, address, name),
              servicesStore.refreshMetadata(ServiceNames.Systems, address, name)
            ]);
          })
          .catch(e => {
            console.error(e);
          });
      return {
        ...node
      };
    });
  });

  /**
   * Check if the node has a proxy system service configured
   *
   * @param {string} nodeName
   * @return {boolean}
   */
  function isProxy(nodeName) {
    return nodeDetails[nodeName]?.systems.metadataTracker?.response?.typeCountsMap?.some(
        ([name, count]) => name === 'proxy' && count > 0
    );
  }
  /**
   * Check if the node has a hub system service configured
   *
   * @param {string} nodeName
   * @return {boolean}
   */
  function isHub(nodeName) {
    return nodeDetails[nodeName]?.systems.metadataTracker?.response?.typeCountsMap?.some(
        ([name, count]) => name === 'hub' && count > 0
    );
  }

  // Clean up on unmount
  onUnmounted(() => {
    unwatchTrackers = [];
    closeResource(hubNodeValue);
  });

  return {
    hubNodeValue,
    enrollHubNodeValue,
    enrollHubNodeAction,

    nodeDetails,
    nodesList,
    isProxy,
    isHub

  };
}
