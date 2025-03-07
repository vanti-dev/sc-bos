import {clientOptions} from '@/api/grpcweb';
import {trackAction} from '@/api/resource';
import {EnrollmentApiPromiseClient} from '@vanti-dev/sc-bos-ui-gen/proto/enrollment_grpc_web_pb';
import {GetEnrollmentRequest, TestEnrollmentRequest} from '@vanti-dev/sc-bos-ui-gen/proto/enrollment_pb';

/**
 *
 * @param {ActionTracker<GetEnrollmentResponse.AsObject>} [tracker]
 * @return {Promise<GetEnrollmentResponse.AsObject>}
 */
export function getEnrollment(tracker) {
  return trackAction('Enrollment.getEnrollment', tracker ?? {}, endpoint => {
    const api = apiClient(endpoint);
    return api.getEnrollment(new GetEnrollmentRequest());
  });
}

/**
 * @param {ActionTracker<TestEnrollmentResponse.AsObject>} [tracker]
 * @return {Promise<GetEnrollmentResponse.AsObject>}
 */
export function testEnrollment(tracker) {
  return trackAction('Enrollment.testEnrollment', tracker ?? {}, endpoint => {
    const api = apiClient(endpoint);
    return api.testEnrollment(new TestEnrollmentRequest());
  });
}

/**
 * @param {string} endpoint
 * @return {EnrollmentApiPromiseClient}
 */
function apiClient(endpoint) {
  return new EnrollmentApiPromiseClient(endpoint, null, clientOptions());
}

