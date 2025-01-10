import {fieldMaskFromObject, setProperties} from '@/api/convpb';
import {clientOptions} from '@/api/grpcweb.js';
import {trackAction} from '@/api/resource';
import {pullResource, setValue} from '@/api/resource.js';
import {periodFromObject} from '@/api/sc/types/period';
import {OccupancySensorHistoryPromiseClient} from '@vanti-dev/sc-bos-ui-gen/proto/history_grpc_web_pb';
import {ListOccupancyHistoryRequest} from '@vanti-dev/sc-bos-ui-gen/proto/history_pb';
import {OccupancySensorApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/occupancy_sensor_grpc_web_pb';
import {Occupancy, PullOccupancyRequest} from '@smart-core-os/sc-api-grpc-web/traits/occupancy_sensor_pb';

/**
 *
 * @param {Partial<PullOccupancyRequest.AsObject>} request
 * @param {ResourceValue<Occupancy.AsObject, PullOccupancyResponse>} resource
 */
export function pullOccupancy(request, resource) {
  pullResource('OccupancySensor.pullOccupancy', resource, endpoint => {
    const api = apiClient(endpoint);
    const stream = api.pullOccupancy(pullOccupancyRequestFromObject(request));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getOccupancy().toObject());
      }
    });
    return stream;
  });
}

/**
 *
 * @param {Partial<ListOccupancyHistoryRequest.AsObject>} request
 * @param {ResourceValue<ListOccupancyHistoryResponse.AsObject>} tracker
 * @return {Promise<ListOccupancyHistoryResponse.AsObject>}
 */
export function listOccupancySensorHistory(request, tracker) {
  return trackAction('OccupancySensorHistory.listOccupancySensorHistory', tracker, (endpoint) => {
    const api = historyClient(endpoint);
    return api.listOccupancyHistory(listOccupancySensorHistoryRequestFromObject(request));
  });
}

/**
 * @param {string} endpoint
 * @return {OccupancySensorApiPromiseClient}
 */
function apiClient(endpoint) {
  return new OccupancySensorApiPromiseClient(endpoint, null, clientOptions());
}

/**
 *
 * @param {string} endpoint
 * @return {OccupancySensorHistoryPromiseClient}
 */
function historyClient(endpoint) {
  return new OccupancySensorHistoryPromiseClient(endpoint, null, clientOptions());
}

/**
 * @param {Partial<PullOccupancyRequest.AsObject>} obj
 * @return {undefined|PullOccupancyRequest}
 */
function pullOccupancyRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new PullOccupancyRequest();
  setProperties(dst, obj, 'name', 'updatesOnly');
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  return dst;
}

/**
 * @param {Partial<ListOccupancyHistoryRequest.AsObject>} obj
 * @return {ListOccupancyHistoryRequest|undefined}
 */
function listOccupancySensorHistoryRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new ListOccupancyHistoryRequest();
  setProperties(dst, obj, 'name', 'pageToken', 'pageSize');
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  dst.setPeriod(periodFromObject(obj.period));
  return dst;
}

/**
 *
 * @param {Occupancy.State} state
 * @return {string}
 */
export function occupancyStateToString(state) {
  switch (state) {
    case Occupancy.State.STATE_UNSPECIFIED:
      return 'Unspecified';
    case Occupancy.State.OCCUPIED:
      return 'Occupied';
    case Occupancy.State.UNOCCUPIED:
      return 'Unoccupied';
    case Occupancy.State.IDLE:
      return 'Idle';
    default:
      return 'Unknown';
  }
}
