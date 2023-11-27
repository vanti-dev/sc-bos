import {clientOptions} from '@/api/grpcweb';
import {setProperties} from '@/api/convpb';
import {pullResource, setCollection, trackAction} from '@/api/resource';
import {HubApiPromiseClient} from '@sc-bos/ui-gen/proto/hub_grpc_web_pb';
import {
  EnrollHubNodeRequest,
  ForgetHubNodeRequest,
  HubNode,
  InspectHubNodeRequest,
  ListHubNodesRequest,
  PullHubNodesRequest,
  TestHubNodeRequest
} from '@sc-bos/ui-gen/proto/hub_pb';

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
 * @param {ActionTracker<HubNode.AsObject>} [tracker]
 * @return {Promise<HubNode.AsObject>}
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
 * @param {ForgetHubNodeRequest.AsObject} request
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
 * @param {InspectHubNodeRequest.AsObject} request
 * @param {ActionTracker<HubNode.AsObject>} [tracker]
 * @return {Promise<HubNode.AsObject>}
 */
export function inspectHubNode(request, tracker) {
  return trackAction('Hub.inspectHubNode', tracker ?? {}, endpoint => {
    const api = apiClient(endpoint);
    return api.inspectHubNode(inspectHubNodeRequestFromObject(request));
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
// ----- Test Hub Node ----- //
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

/**
 *
 * @param {InspectHubNodeRequest.AsObject} obj
 * @return {InspectHubNodeRequest|undefined}
 */
function inspectHubNodeRequestFromObject(obj) {
  if (!obj) return undefined;

  const dst = new InspectHubNodeRequest();
  dst.setNode(hubNodeFromObject(obj.node));
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
  setProperties(dst, obj, 'publicCertsList');
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
function forgetHubNodeRequestFromObject(obj) {
  if (!obj) return undefined;

  const dst = new ForgetHubNodeRequest();
  setProperties(dst, obj, 'address', 'allowMissing');

  return dst;
}

