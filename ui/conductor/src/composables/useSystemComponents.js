import {computed, onUnmounted, reactive, ref, set, watch, watchEffect} from 'vue';
import {closeResource, newActionTracker} from '@/api/resource';
import {enrollHubNode, forgetHubNode, inspectHubNode, testHubNode} from '@/api/sc/traits/hub';
import {ServiceNames} from '@/api/ui/services';
import {useEnrollmentStore} from '@/stores/enrollmentStore';
import {useHubStore} from '@/stores/hub';
import {useServicesStore} from '@/stores/services';
import {parseCertificate} from '@/util/certificates';

// generate ServiceTracker type
/**
 * @typedef {import('@/api/ui/services').ServiceTracker} ServiceTracker
 * @return {{
 *   hubNodeValue: ActionTracker<TestHubNodeResponse.AsObject>,
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
  const forgetHubNodeValue = reactive(newActionTracker());
  const inspectHubNodeValue = reactive(newActionTracker());

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
      publicCertsList: inspectHubNodeValue.response.publicCertsList
    };

    return enrollHubNode(request, enrollHubNodeValue);
  }

  /**
   *
   * @param {string} address
   * @return {Promise<void>|undefined}
   */
  function forgetHubNodeAction(address) {
    if (!address) return;

    const request = {
      address,
      allowMissing: true
    };

    return forgetHubNode(request, forgetHubNodeValue);
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
    const {traitsList, ...metadata} = inspectHubNodeValue.response.metadata;

    return metadata;
  });
  // --------------------------- //
  // List and Track Hub Nodes
  const nodeDetails = reactive({});

  let unwatchTrackers = [];

  const nodesList = computed(() => {
    return Object.values(hubStore.nodesList).map(node => {
      Promise.all([node.commsAddress, node.commsName])
          .then(([address, name]) => {
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

  const allowForget = (name) => {
    return !!hubStore.listedHubNodes.find(node => node === name);
  };

  // Clean up on unmount
  onUnmounted(() => {
    unwatchTrackers = [];
    closeResource(hubNodeValue);
  });

  return {
    hubNodeValue,
    enrollHubNodeValue,
    forgetHubNodeValue,
    inspectHubNodeValue,
    enrollHubNodeAction,
    forgetHubNodeAction,
    inspectHubNodeAction,
    readCertificates,
    resetCertificates,
    readMetadata,

    nodeDetails,
    nodesList,
    isProxy,
    isHub,
    allowForget

  };
}
