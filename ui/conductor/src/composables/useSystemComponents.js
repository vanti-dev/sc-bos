import {closeResource, newActionTracker} from '@/api/resource';
import {enrollHubNode, forgetHubNode, inspectHubNode, testHubNode} from '@/api/sc/traits/hub';
import {ServiceNames} from '@/api/ui/services';
import {useHubStore} from '@/stores/hub';
import {useServicesStore} from '@/stores/services';
import {parseCertificate} from '@/util/certificates';
import {computed, onUnmounted, reactive, ref, watchEffect} from 'vue';

/**
 * @typedef {import('@/api/ui/services').ServiceTracker} ServiceTracker
 * @return {{
 *   enrollHubNodeValue: ActionTracker<HubNode.AsObject>,
 *   enrollHubNodeAction: (address: string) => Promise<void>,
 *   forgetHubNodeAction: (address: string) => Promise<void>,
 *   forgetHubNodeValue: ActionTracker<HubNode.AsObject>,
 *   inspectHubNodeValue: ActionTracker<InspectHubNodeResponse.AsObject>,
 *   inspectHubNodeAction: (address: string) => Promise<void>,
 *   readCertificates: import('vue').ComputedRef<{
 *     validityPeriod: string,
 *     extensions: {
 *       keyUsage: string,
 *       basicConstraints: string,
 *       subjectKeyIdentifier: string
 *     },
 *     keyLength: number,
 *     serial: string,
 *     subject: {
 *       commonName: string,
 *       organization: string
 *     },
 *     sha1Fingerprint: string,
 *     sha256Fingerprint: string,
 *     primaryDomain: string,
 *     subjectAltDomains: string,
 *     version: number,
 *     signatureAlgorithm: string,
 *     issuer: {
 *       commonName: string
 *     }
 *   }[]>,
 *   resetCertificates: () => void,
 *   readMetadata: import('vue').ComputedRef<MetadataResponse.AsObject>,
 *   nodeDetails: Record<string, {
 *     automations: ServiceTracker,
 *     drivers: ServiceTracker,
 *     systems: ServiceTracker
 *   }>,
 *   nodesList: import('vue').ComputedRef<HubNode.AsObject[]>,
 *   isProxy: (nodeName: string) => boolean,
 *   isHub: (nodeName: string) => boolean
 * }}
 */
export default function() {
  const hubStore = useHubStore();
  const servicesStore = useServicesStore();


  // --------------------------- //
  // Manage Hub Nodes
  const testHubNodeValue = reactive(
      /** @type {ActionTracker<TestHubNodeResponse.AsObject>} */ newActionTracker()
  );
  const enrollHubNodeValue = reactive(
      /** @type {ActionTracker<HubNode.AsObject>} */ newActionTracker()
  );
  const forgetHubNodeValue = reactive(
      /** @type {ActionTracker<ForgetHubNodeResponse.AsObject>} */ newActionTracker()
  );
  const inspectHubNodeValue = reactive(
      /** @type {ActionTracker<HubNode.AsObject>} */ newActionTracker()
  );

  /**
   *
   * @param {string} address
   * @return {Promise<void>|undefined}
   */
  async function testHubNodeAction(address) {
    if (!address) return;

    const request = {
      address
    };

    await testHubNode(request, testHubNodeValue);
  }

  /**
   * @param {string} address
   * @return {Promise<void>|undefined}
   */
  async function enrollHubNodeAction(address) {
    if (!address) return;

    const request = {
      node: {
        address
      },
      publicCertsList: inspectHubNodeValue.response.publicCertsList
    };

    // Enroll the node
    await enrollHubNode(request, enrollHubNodeValue);

    // Refresh the list of nodes
    await hubStore.listHubNodesAction();
  }

  /**
   *
   * @param {string} address
   * @return {Promise<void>|undefined}
   */
  async function forgetHubNodeAction(address) {
    if (!address) return;

    const request = {
      address,
      allowMissing: true
    };

    await forgetHubNode(request, forgetHubNodeValue);

    await hubStore.listHubNodesAction();
  }

  /**
   *
   * @param {string} address
   * @return {Promise<void>|undefined}
   */
  function inspectHubNodeAction(address) {
    if (!address) return;

    const request = {
      node: {
        address
      }
    };

    return inspectHubNode(request, inspectHubNodeValue);
  }

  const resetInspectHubNodeValue = () => {
    Object.assign(inspectHubNodeValue, newActionTracker());
  };

  const parsedCertificatesData = ref([]);

  // Computed property to use in the template
  const readCertificates = computed(() => parsedCertificatesData.value);

  // Watcher to react to changes in inspectHubNodeValue
  watchEffect(async () => {
    if (!inspectHubNodeValue?.response?.publicCertsList) {
      parsedCertificatesData.value = [];
      return;
    }

    const parsedCertificatesPromises = inspectHubNodeValue.response.publicCertsList.map(cert => parseCertificate(cert));
    parsedCertificatesData.value = await Promise.all(parsedCertificatesPromises);
  });

  // Removing previously read node details
  const resetCertificates = () => {
    parsedCertificatesData.value = [];
    Object.assign(inspectHubNodeValue, newActionTracker());
  };

  const readMetadata = computed(() => {
    if (!inspectHubNodeValue?.response?.metadata) {
      return null;
    }

    // exclude traitsList from the metadata return
    // eslint-disable-next-line no-unused-vars
    const {traitsList, ...metadata} = inspectHubNodeValue.response.metadata;

    return metadata;
  });
  // --------------------------- //
  // List and Track Hub Nodes
  const nodeDetails = reactive({});
  const nodesList = ref([]);
  let unwatchTrackers = [];

  const processNodes = () => {
    const nodes = Object.values(hubStore.nodesList);
    nodesList.value = [];

    nodesList.value = nodes.map(node => {
      // Using an immediately-invoked async function to handle asynchronous operations
      (async () => {
        try {
          const [address, name] = await Promise.all([node.commsAddress, node.commsName]);

          nodeDetails[node.name] = {
            automations: servicesStore.getService(ServiceNames.Automations, address, name),
            drivers: servicesStore.getService(ServiceNames.Drivers, address, name),
            systems: servicesStore.getService(ServiceNames.Systems, address, name)
          };

          unwatchTrackers.push(nodeDetails[node.name].automations.metadataTracker);
          unwatchTrackers.push(nodeDetails[node.name].drivers.metadataTracker);
          unwatchTrackers.push(nodeDetails[node.name].systems.metadataTracker);

          await Promise.all([
            servicesStore.refreshMetadata(ServiceNames.Automations, address, name),
            servicesStore.refreshMetadata(ServiceNames.Drivers, address, name),
            servicesStore.refreshMetadata(ServiceNames.Systems, address, name)
          ]);
        } catch (e) {
          console.error('Error processing node:', e);
        }
      })();

      return {
        ...node
      };
    });
  };

  watchEffect(() => {
    processNodes();
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

  const allowForget = (name) => {
    return !!hubStore.listedHubNodes.find(node => node === name);
  };

  // Clean up on unmount
  onUnmounted(() => {
    unwatchTrackers = [];
    closeResource(enrollHubNodeValue);
    closeResource(forgetHubNodeValue);
    closeResource(inspectHubNodeValue);
    closeResource(testHubNodeValue);
  });

  return {
    testHubNodeValue,
    enrollHubNodeValue,
    forgetHubNodeValue,
    inspectHubNodeValue,
    testHubNodeAction,
    enrollHubNodeAction,
    forgetHubNodeAction,
    inspectHubNodeAction,
    resetInspectHubNodeValue,
    readCertificates,
    resetCertificates,
    readMetadata,

    nodeDetails,
    nodesList,
    processNodes,
    isProxy,
    isHub,
    allowForget

  };
}
