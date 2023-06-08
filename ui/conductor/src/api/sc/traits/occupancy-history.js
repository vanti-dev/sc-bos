import {clientOptions} from '@/api/grpcweb';
import {trackAction} from '@/api/resource';
import {OccupancySensorHistoryPromiseClient} from '@sc-bos/ui-gen/proto/history_grpc_web_pb';
import {ListOccupancyHistoryRequest} from '@sc-bos/ui-gen/proto/history_pb';
import {fieldMaskFromObject, setProperties} from '@/api/convpb';
import {periodFromObject} from '@/api/sc/types/period';

/**
 *
 * @param {ListOccupancyHistoryRequest.AsObject} request
 * @param {ResourceValue<ListOccupancyHistoryResponse.AsObject>} tracker
 * @return {Promise<ListOccupancyHistoryResponse.AsObject>}
 */
export function listOccupancySensorHistory(request, tracker) {
  return trackAction('OccupancySensorHistory.listOccupancySensorHistory', tracker, (endpoint) => {
    const api = client(endpoint);
    return api.listOccupancyHistory(listOccupancySensorHistoryRequestFromObject(request));
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

/**
 * @param {ListOccupancyHistoryRequest.AsObject} obj
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
