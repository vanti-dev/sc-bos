import {OnOffApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/on_off_grpc_web_pb.js';
import {GetOnOffRequest, PullOnOffRequest} from '@smart-core-os/sc-api-grpc-web/traits/on_off_pb.js';
import {pullResource, setValue, trackAction} from './resource.js';
import {clientOptions} from '../../grpcweb.js';

/**
 * @param {string} name
 * @param {ResourceValue<OnOff.AsObject, OnOff>} resource
 */
export function pullOnOff(name, resource) {
  pullResource('OnOff', resource, endpoint => {
    const api = new OnOffApiPromiseClient(endpoint, null, clientOptions());
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

/**
 * @param {string} name
 * @param {ActionTracker<OnOff.AsObject>} [tracker]
 * @return {Promise<OnOff.AsObject>}
 */
export function getOnOff(name, tracker) {
  return trackAction('OnOff.createOnOff', tracker ?? {}, endpoint => {
    const api = new OnOffApiPromiseClient(endpoint, null, clientOptions());
    return api.getOnOff(
        new GetOnOffRequest()
            .setName(name)
    );
  });
}
