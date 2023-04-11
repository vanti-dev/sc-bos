import {clientOptions} from '@/api/grpcweb.js';
import {pullResource, setError, setValue} from '@/api/resource.js';
import {OpenCloseApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/open_close_grpc_web_pb';
import {PullOpenClosePositionsRequest} from '@smart-core-os/sc-api-grpc-web/traits/open_close_pb';

/**
 * @param {string} name
 * @param {ResourceValue<OpenClosePositions.AsObject, OpenClosePositions>} resource
 */
export function pullOpenClosePositions(name, resource) {
  pullResource('OpenClose.Positions', resource, endpoint => {
    const api = new OpenCloseApiPromiseClient(endpoint, null, clientOptions());
    const stream = api.pullPositions(new PullOpenClosePositionsRequest().setName(name));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getOpenClosePosition().toObject());
      }
    });
    stream.on('error', err => {
      setError(resource, err);
    });
    return stream;
  });
}
