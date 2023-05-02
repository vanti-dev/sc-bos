import {clientOptions} from '@/api/grpcweb';
import {pullResource, setCollection, trackAction} from '@/api/resource';
import {HubApiPromiseClient} from '@sc-bos/ui-gen/proto/hub_grpc_web_pb';
import {ListHubNodesRequest, PullHubNodesRequest} from '@sc-bos/ui-gen/proto/hub_pb';

/**
 *
 * @param {ActionTracker<ListHubNodesResponse.AsObject>} tracker
 * @return {Promise<ListHubNodesResponse.AsObject>}
 */
export function listHubNodes(tracker) {
  return trackAction('Hub.listHubNodes', tracker ?? {}, endpoint => {
    const api = new HubApiPromiseClient(endpoint, null, clientOptions());
    return api.listHubNodes(new ListHubNodesRequest());
  });
}

/**
 *
 * @param {ResourceCollection<HubNode.AsObject, HubNode>} resource
 */
export function pullHubNodes(resource) {
  pullResource('Hub.pullHubNodes', resource, endpoint => {
    const api = new HubApiPromiseClient(endpoint, null, clientOptions());
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
