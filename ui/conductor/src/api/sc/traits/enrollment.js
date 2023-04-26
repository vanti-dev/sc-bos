import {clientOptions} from '@/api/grpcweb';
import {trackAction} from '@/api/resource';
import {EnrollmentApiPromiseClient} from '@sc-bos/ui-gen/proto/enrollment_grpc_web_pb';
import {GetEnrollmentRequest} from '@sc-bos/ui-gen/proto/enrollment_pb';

/**
 *
 * @param {ActionTracker<GetEnrollmentResponse.AsObject>} tracker
 * @return {Promise<GetEnrollmentResponse.AsObject>}
 */
export function getEnrollment(tracker) {
  return trackAction('Enrollment.getEnrollment', tracker ?? {}, endpoint => {
    const api = new EnrollmentApiPromiseClient(endpoint, null, clientOptions());
    return api.getEnrollment(new GetEnrollmentRequest());
  });
}

