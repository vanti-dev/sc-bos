import {newActionTracker} from '@/api/resource';
import {enrollHubNode, forgetHubNode, inspectHubNode} from '@/api/sc/traits/hub';
import {parseCertificate} from '@/util/certificates';
import {computed, reactive, ref, watchEffect} from 'vue';

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
 *   readMetadata: import('vue').ComputedRef<MetadataResponse.AsObject>
 * }}
 */
export default function() {
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

  return {
    testHubNodeValue,
    enrollHubNodeValue,
    forgetHubNodeValue,
    inspectHubNodeValue,
    enrollHubNodeAction,
    forgetHubNodeAction,
    inspectHubNodeAction,
    resetInspectHubNodeValue,
    readCertificates,
    resetCertificates,
    readMetadata
  };
}
