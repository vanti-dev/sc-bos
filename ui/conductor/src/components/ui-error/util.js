import {StatusCode} from 'grpc-web';

/**
 * @param {StatusCode} statusCode
 * @return {string}
 */
export function statusCodeToString(statusCode) {
  switch (statusCode) {
    case StatusCode.OK:
      return 'OK';
    case StatusCode.CANCELLED:
      return 'CANCELLED';
    case StatusCode.UNKNOWN:
      return 'UNKNOWN';
    case StatusCode.INVALID_ARGUMENT:
      return 'INVALID ARGUMENT';
    case StatusCode.DEADLINE_EXCEEDED:
      return 'DEADLINE EXCEEDED';
    case StatusCode.NOT_FOUND:
      return 'NOT FOUND';
    case StatusCode.ALREADY_EXISTS:
      return 'ALREADY EXISTS';
    case StatusCode.PERMISSION_DENIED:
      return 'PERMISSION DENIED';
    case StatusCode.RESOURCE_EXHAUSTED:
      return 'RESOURCE EXHAUSTED';
    case StatusCode.FAILED_PRECONDITION:
      return 'FAILED PRECONDITION';
    case StatusCode.ABORTED:
      return 'ABORTED';
    case StatusCode.OUT_OF_RANGE:
      return 'OUT OF RANGE';
    case StatusCode.UNIMPLEMENTED:
      return 'UNIMPLEMENTED';
    case StatusCode.INTERNAL:
      return 'INTERNAL';
    case StatusCode.UNAVAILABLE:
      return 'UNAVAILABLE';
    case StatusCode.DATA_LOSS:
      return 'DATA LOSS';
    case StatusCode.UNAUTHENTICATED:
      return 'UNAUTHENTICATED';
    default:
      return 'UNKNOWN ERROR';
  }
}
