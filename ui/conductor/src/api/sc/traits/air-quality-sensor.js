import {clientOptions} from '@/api/grpcweb.js';
import {pullResource, setValue} from '@/api/resource.js';
import {AirQualitySensorApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/air_quality_sensor_grpc_web_pb';
import {PullAirQualityRequest} from '@smart-core-os/sc-api-grpc-web/traits/air_quality_sensor_pb';

/**
 * @param {string} name
 * @param {ResourceValue<AirQuality.AsObject, AirQuality>} resource
 */
export function pullAirQualitySensor(name, resource) {
  pullResource('AirQualitySensor', resource, endpoint => {
    const api = new AirQualitySensorApiPromiseClient(endpoint, null, clientOptions());
    const stream = api.pullAirQuality(new PullAirQualityRequest().setName(name));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getAirQuality().toObject());
      }
    });
    return stream;
  });
}
