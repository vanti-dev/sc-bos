import {LightingTestApiPromiseClient} from '@sc-bos/ui-gen/proto/lighting_test_grpc_web_pb';
import {GetReportCSVRequest, ListLightHealthRequest} from '@sc-bos/ui-gen/proto/lighting_test_pb';
import {clientOptions} from '@/api/grpcweb';
import {trackAction} from '@/api/resource';

/**
 *
 * @param {ActionTracker<ReportCSV.AsObject>} tracker
 */
export function getReportCSV(tracker) {
  return trackAction('LightingTest.getReportCSV', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    const req = new GetReportCSVRequest();
    req.setIncludeHeader(true);
    return api.getReportCSV(req);
  });
}

/**
 *
 * @param {ActionTracker<ListLightHealthResponse.AsObject>} tracker
 */
export function listLightHealth(tracker) {
  return trackAction('LightingTest.listLightHealth', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.listLightHealth(new ListLightHealthRequest());
  });
}

/**
 *
 * @param {number} faultId
 * @return {string}
 */
export function faultToString(faultId) {
  switch (faultId) {
    case 1: return 'DURATION_TEST_FAILED';
    case 2: return 'FUNCTION_TEST_FAILED';
    case 3: return 'BATTERY_FAULT';
    case 4: return 'LAMP_FAULT';
    case 5: return 'COMMUNICATION_FAILURE';
    case 6: return 'OTHER_FAULT';
    default: return 'FAULT_UNSPECIFIED';
  }
}

/**
 *
 * @param {string} endpoint
 */
function client(endpoint) {
  return new LightingTestApiPromiseClient(endpoint, null, clientOptions());
}
