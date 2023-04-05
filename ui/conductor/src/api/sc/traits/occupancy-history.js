import {clientOptions} from '@/api/grpcweb';
import {trackAction} from '@/api/resource';
import {
  MeterHistoryPromiseClient,
  OccupancySensorHistoryClient,
  OccupancySensorHistoryPromiseClient
} from '@sc-bos/ui-gen/proto/history_grpc_web_pb';
import {ListMeterReadingHistoryRequest, ListOccupancyHistoryRequest} from '@sc-bos/ui-gen/proto/history_pb';

/**
 *
 * @param {string} name
 * @param {ResourceValue<ListOccupancyHistoryResponse.AsObject>} tracker
 * @return {Promise<ListOccupancyHistoryResponse.AsObject>}
 */
export function listOccupancySensorHistory(name, tracker) {
  return trackAction('OccupancySensorHistory.listOccupancySensorHistory', tracker, endpoint => {
    const api = client(endpoint);
    return api.listOccupancyHistory(new ListOccupancyHistoryRequest().setName(name));
  });
}


/**
 *
 * @param {string} endpoint
 * @return {OccupancySensorHistoryPromiseClient}
 */
function client(endpoint) {
  return new OccupancySensorHistoryPromiseClient(endpoint, null, clientOptions());
}
