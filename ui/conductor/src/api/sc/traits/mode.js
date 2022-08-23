import {clientOptions} from '@/api/grpcweb.js';
import {pullResource, setValue} from '@/api/resource.js';
import {ModeApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/mode_grpc_web_pb.js';
import {PullModeValuesRequest} from '@smart-core-os/sc-api-grpc-web/traits/mode_pb.js';

/**
 * @param {string} name
 * @param {ResourceValue<ModeValues.AsObject, ModeValues>} resource
 */
export function pullModeValues(name, resource) {
  pullResource('Mode.ModeValues', resource, endpoint => {
    const api = new ModeApiPromiseClient(endpoint, null, clientOptions());
    const stream = api.pullModeValues(new PullModeValuesRequest().setName(name));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getModeValues().toObject());
      }
    });
    return stream;
  });
}
