import {fieldMaskFromObject, setProperties} from '@/api/convpb';
import {clientOptions} from '@/api/grpcweb.js';
import {pullResource, setValue} from '@/api/resource.js';
import {AirQualitySensorApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/air_quality_sensor_grpc_web_pb';
import {PullAirQualityRequest} from '@smart-core-os/sc-api-grpc-web/traits/air_quality_sensor_pb';

/**
 * @param {PullAirQualityRequest.AsObject} request
 * @param {ResourceValue<AirQuality.AsObject, PullAirQualityResponse>} resource
 */
export function pullAirQualitySensor(request, resource) {
  pullResource('AirQualitySensor.pullAirQuality', resource, endpoint => {
    const api = apiClient(endpoint);
    const stream = api.pullAirQuality(pullAirQualityRequestFromObject(request));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getAirQuality().toObject());
      }
    });
    return stream;
  });
}

/**
 * @param {string} endpoint
 * @return {AirQualitySensorApiPromiseClient}
 */
function apiClient(endpoint) {
  return new AirQualitySensorApiPromiseClient(endpoint, null, clientOptions());
}

/**
 * @param {PullAirQualityRequest.AsObject} obj
 * @return {PullAirQualityRequest|undefined}
 */
function pullAirQualityRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new PullAirQualityRequest();
  setProperties(dst, obj, 'name', 'updatesOnly');
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  return dst;
}
