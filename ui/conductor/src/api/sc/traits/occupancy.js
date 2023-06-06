import {clientOptions} from '@/api/grpcweb.js';
import {pullResource, setValue} from '@/api/resource.js';
import {OccupancySensorApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/occupancy_sensor_grpc_web_pb';
import {Occupancy, PullOccupancyRequest} from '@smart-core-os/sc-api-grpc-web/traits/occupancy_sensor_pb';

/**
 *
 * @param {string} name
 * @param {ResourceValue<Occupancy.AsObject, Occupancy>} resource
 */
export function pullOccupancy(name, resource) {
  pullResource('OccupancySensor.Occupancy', resource, endpoint => {
    // console.log(name);
    const api = new OccupancySensorApiPromiseClient(endpoint, null, clientOptions());
    const stream = api.pullOccupancy(new PullOccupancyRequest().setName(name));
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
