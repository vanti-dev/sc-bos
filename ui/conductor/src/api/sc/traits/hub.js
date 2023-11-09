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
 * @param {string} address
 * @param {ActionTracker<EnrollHubNodeResponse.AsObject>} [tracker]
 * @return {Promise<EnrollHubNodeResponse.AsObject>}
 */
export function enrollHubNode(address, tracker) {
  return trackAction('Hub.enrollHubNode', tracker ?? {}, endpoint => {
    const api = apiClient(endpoint);
    const hubNode = setHubNodeFromObject(enrollHubNodeFromObject(address));
    return api.enrollHubNode(hubNode);
  });
}

/**
 *
 * @param {string} address
 * @param {ActionTracker<ForgetHubNodeResponse.AsObject>} [tracker]
 * @return {Promise<ForgetHubNodeResponse.AsObject>}
 */
export function forgetHubNode(address, tracker) {
  return trackAction('Hub.forgetHubNode', tracker ?? {}, endpoint => {
    const api = apiClient(endpoint);
    const request = new ForgetHubNodeRequest();
    request.setAddress(address);
    request.setAllowMissing(true);
    return api.forgetHubNode(request);
  });
}

/**
 * @param {string} address
 * @param {ActionTracker<TestHubNodeResponse.AsObject>} [tracker]
 * @return {Promise<TestHubNodeResponse.AsObject>}
 */
export function testHubNodes(address, tracker) {
  return trackAction('Hub.testHubNode', tracker ?? {}, endpoint => {
    const api = apiClient(endpoint);
    return api.testHubNode(testHubNodesFromObject({address}));
  });
}

/**
 * @param {string} endpoint
 * @return {HubApiPromiseClient}
 */
function apiClient(endpoint) {
  return new HubApiPromiseClient(endpoint, null, clientOptions());
}

/**
 *
 * @param {TestHubNodeRequest.AsObject} obj
 * @return {TestHubNodeRequest|undefined}
 */
function testHubNodesFromObject(obj) {
  if (!obj) return undefined;

  const dst = new TestHubNodeRequest();
  dst.setAddress(obj.address);
  return dst;
}

/**
 * @param {HubNode.AsObject} obj
 * @return {EnrollHubNodeRequest|undefined}
 */
function setHubNodeFromObject(obj) {
  if (!obj) return undefined;

  const dst = new EnrollHubNodeRequest();
  dst.setNode(obj);

  return dst;
}


/**
 * @param {HubNode.AsObject} obj
 * @return {HubNode|undefined}
 */
function enrollHubNodeFromObject(obj) {
  if (!obj) return undefined;

  const dst = new HubNode();
  dst.setAddress(obj);

  return dst;
}
