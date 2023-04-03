import {clientOptions} from '@/api/grpcweb';
import {trackAction} from '@/api/resource';
import {MeterHistoryPromiseClient} from '@sc-bos/ui-gen/proto/history_grpc_web_pb';
import {ListMeterReadingHistoryRequest} from '@sc-bos/ui-gen/proto/history_pb';

/**
 *
 * @param {string} name
 * @param {ResourceValue<ListMeterReadingHistoryResponse.AsObject>} tracker
 * @return {Promise<ListMeterReadingHistoryResponse.AsObject>}
 */
export function listMeterReadingHistory(name, tracker) {
  return trackAction('MeterReadingHistory.listMeterReadingHistory', tracker, endpoint => {
    const api = client(endpoint);
    return api.listMeterReadingHistory(new ListMeterReadingHistoryRequest().setName(name));
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
