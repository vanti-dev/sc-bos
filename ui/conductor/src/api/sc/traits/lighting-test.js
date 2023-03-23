import {LightingTestApiPromiseClient} from '@sc-bos/ui-gen/proto/lighting_test_grpc_web_pb';
import {clientOptions} from '@/api/grpcweb';
import {GetReportCSVRequest} from '@sc-bos/ui-gen/proto/lighting_test_pb';
import {trackAction} from '@/api/resource';

/**
 *
 * @param {ActionTracker<ReportCSV.AsObject>} tracker
 */
export function getReportCSV(tracker) {
  return trackAction('LightingTest.GetReportCSV', tracker ?? {}, endpoint => {
    const api = new LightingTestApiPromiseClient(endpoint, null, clientOptions());
    const req = new GetReportCSVRequest();
    req.setIncludeHeader(true);
    return api.getReportCSV(req);
  });
}
