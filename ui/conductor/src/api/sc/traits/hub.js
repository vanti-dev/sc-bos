import {clientOptions} from '@/api/grpcweb';
import {pullResource, setCollection, trackAction} from '@/api/resource';
import {HubApiPromiseClient} from '@sc-bos/ui-gen/proto/hub_grpc_web_pb';
import {
  EnrollHubNodeRequest,
  ForgetHubNodeRequest,
  HubNode,
  ListHubNodesRequest,
  PullHubNodesRequest,
  TestHubNodeRequest
} from '@sc-bos/ui-gen/proto/hub_pb';
import {setProperties} from "@/api/convpb";

/**
 *
 * @param {ActionTracker<ListHubNodesResponse.AsObject>} [tracker]
 * @return {Promise<ListHubNodesResponse.AsObject>}
 */
export function listHubNodes(tracker) {
  return trackAction('Hub.listHubNodes', tracker ?? {}, endpoint => {
    const api = apiClient(endpoint);
    return api.listHubNodes(new ListHubNodesRequest());
  });
}

/**
 *
 * @param {ResourceCollection<HubNode.AsObject, PullHubNodesResponse>} resource
 */
export function pullHubNodes(resource) {
  pullResource('Hub.pullHubNodes', resource, endpoint => {
    const api = apiClient(endpoint);
    const stream = api.pullHubNodes(new PullHubNodesRequest());
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setCollection(resource, change, v => v.name);
      }
    });
    return stream;
  });
}


/**
 *
 * @param {EnrollHubNodeRequest.AsObject} request
 * @param {ActionTracker<EnrollHubNodeResponse.AsObject>} [tracker]
 * @return {Promise<EnrollHubNodeResponse.AsObject>}
 */
export function enrollHubNode(request, tracker) {
  return trackAction('Hub.enrollHubNode', tracker ?? {}, endpoint => {
    const api = apiClient(endpoint);
    const hubNode = enrollHubNodeRequestFromObject(request);
    return api.enrollHubNode(hubNode);
  });
}

/**
 *
 * @param {HubNode.AsObject} request
 * @param {ActionTracker<ForgetHubNodeResponse.AsObject>} [tracker]
 * @return {Promise<ForgetHubNodeResponse.AsObject>}
 */
export function forgetHubNode(request, tracker) {
  return trackAction('Hub.forgetHubNode', tracker ?? {}, endpoint => {
    const api = apiClient(endpoint);
    const hubNode = forgetHubNodeRequestFromObject(request);
    return api.forgetHubNode(hubNode);
  });
}

/**
 * @param {TestHubNodeRequest.AsObject} request
 * @param {ActionTracker<TestHubNodeResponse.AsObject>} [tracker]
 * @return {Promise<TestHubNodeResponse.AsObject>}
 */
export function testHubNode(request, tracker) {
  return trackAction('Hub.testHubNode', tracker ?? {}, endpoint => {
    const api = apiClient(endpoint);
    return api.testHubNode(testHubNodeRequestFromObject(request));
  });
}

/**
 * @param {string} endpoint
 * @return {HubApiPromiseClient}
 */
function apiClient(endpoint) {
  return new HubApiPromiseClient(endpoint, null, clientOptions());
}

// --------------------------- //
// ----- Enroll Hub Node ----- //
/**
 *
 * @param {TestHubNodeRequest.AsObject} obj
 * @return {TestHubNodeRequest|undefined}
 */
function testHubNodeRequestFromObject(obj) {
  if (!obj) return undefined;

  const dst = new TestHubNodeRequest();
  dst.setAddress(obj.address);
  return dst;
}
// --------------------------- //

// --------------------------- //
// ----- Enroll Hub Node ----- //
/**
 * @param {EnrollHubNodeRequest.AsObject} obj
 * @return {HubNode.AsObject|undefined}
 */
function enrollHubNodeRequestFromObject(obj) {
  if (!obj) return undefined;

  const dst = new EnrollHubNodeRequest();
  setProperties(dst, obj, 'node', 'publicCertsList');
  dst.setNode(hubNodeFromObject(obj.node));
  return dst;
}


/**
 * @param {EnrollHubNodeRequest.AsObject} obj
 * @return {HubNode.AsObject|undefined}
 */
function hubNodeFromObject(obj) {
  if (!obj) return undefined;

  const dst = new HubNode();
  setProperties(dst, obj, 'address', 'name', 'description');

  return dst;
}
// --------------------------- //


// --------------------------- //
// ----- Forget Hub Node ----- //
/**
 * @param {HubNode.AsObject} obj
 * @return {ForgetHubNodeRequest|undefined}
 */
function forgetHubNodRequestFromObject(obj) {
  if (!obj) return undefined;

  const dst = new ForgetHubNodeRequest();
  setProperties(dst, obj, 'address', 'allowMissing');

  return dst;
}

