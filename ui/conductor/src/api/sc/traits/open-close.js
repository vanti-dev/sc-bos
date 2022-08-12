import {OpenCloseApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/open_close_grpc_web_pb.js';
import {PullOpenClosePositionsRequest} from '@smart-core-os/sc-api-grpc-web/traits/open_close_pb.js';
import {pullResource, setValue} from './resource.js';

/**
 * @param {string} name
 * @param {ResourceValue<OpenClosePositions.AsObject, OpenClosePositions>} resource
 */
export function pullOpenClosePositions(name, resource) {
  pullResource('OpenClose.Positions', resource, endpoint => {
    const api = new OpenCloseApiPromiseClient(endpoint);
    const stream = api.pullPositions(new PullOpenClosePositionsRequest().setName(name));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getOpenClosePosition().toObject());
      }
    });
    return stream;
  });
}
