import {OnOffApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/on_off_grpc_web_pb.js';
import {PullOnOffRequest} from '@smart-core-os/sc-api-grpc-web/traits/on_off_pb.js';
import {pullResource, setValue} from './resource.js';

/**
 * @param {string} name
 * @param {ResourceValue<OnOff.AsObject, OnOff>} resource
 */
export function pullOnOff(name, resource) {
  pullResource('OnOff', resource, endpoint => {
    const api = new OnOffApiPromiseClient(endpoint);
    const stream = api.pullOnOff(new PullOnOffRequest().setName(name));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getOnOff().toObject());
      }
    });
    return stream;
  });
}
