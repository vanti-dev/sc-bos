import {clientOptions} from '@/api/grpcweb.js';
import {pullResource, setValue} from '@/api/resource.js';
import {LightApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/light_grpc_web_pb.js';
import {PullBrightnessRequest} from '@smart-core-os/sc-api-grpc-web/traits/light_pb.js';

/**
 * @param {string} name
 * @param {ResourceValue<Brightness.AsObject, Brightness>} resource
 */
export function pullBrightness(name, resource) {
  pullResource('Light.Brightness', resource, endpoint => {
    const api = new LightApiPromiseClient(endpoint, null, clientOptions());
    const stream = api.pullBrightness(new PullBrightnessRequest().setName(name));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getBrightness().toObject());
      }
    });
    return stream;
  });
}
