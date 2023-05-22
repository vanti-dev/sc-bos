import {fieldMaskFromObject, setProperties} from '@/api/convpb';
import {clientOptions} from '@/api/grpcweb';
import {trackAction} from '@/api/resource';
import {periodFromObject} from '@/api/sc/types/period';
import {MeterHistoryPromiseClient} from '@sc-bos/ui-gen/proto/history_grpc_web_pb';
import {ListMeterReadingHistoryRequest} from '@sc-bos/ui-gen/proto/history_pb';

/**
 *
 * @param {ListMeterReadingHistoryRequest.AsObject} request
 * @param {ActionTracker<ListMeterReadingHistoryResponse.AsObject>} tracker
 * @return {Promise<ListMeterReadingHistoryResponse.AsObject>}
 */
export function listMeterReadingHistory(request, tracker) {
  return trackAction('MeterReadingHistory.listMeterReadingHistory', tracker, endpoint => {
    const api = client(endpoint);
    return api.listMeterReadingHistory(listMeterReadingHistoryRequestFromObject(request));
  });
}


/**
 *
 * @param {string} endpoint
 * @return {MeterHistoryPromiseClient}
 */
function client(endpoint) {
  return new MeterHistoryPromiseClient(endpoint, null, clientOptions());
}

/**
 * @param {ListMeterReadingHistoryRequest.AsObject} obj
 * @return {ListMeterReadingHistoryRequest|undefined}
 */
function listMeterReadingHistoryRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new ListMeterReadingHistoryRequest();
  setProperties(dst, obj, 'name', 'pageToken', 'pageSize');
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  dst.setPeriod(periodFromObject(obj.period));
  return dst;
}
