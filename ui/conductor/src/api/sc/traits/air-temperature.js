import {AirTemperatureApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/air_temperature_grpc_web_pb.js';
import {PullAirTemperatureRequest} from '@smart-core-os/sc-api-grpc-web/traits/air_temperature_pb.js';
import {pullResource, setValue} from './resource.js';

/**
 * @param {string} name
 * @param {ResourceValue<AirTemperature.AsObject, AirTemperature>} resource
 */
export function pullAirTemperature(name, resource) {
  pullResource('AirTemperature', resource, endpoint => {
    const api = new AirTemperatureApiPromiseClient(endpoint);
    const stream = api.pullAirTemperature(new PullAirTemperatureRequest().setName(name));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getAirTemperature().toObject());
      }
    });
    return stream;
  });
}
