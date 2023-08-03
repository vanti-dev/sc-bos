import {fieldMaskFromObject, setProperties} from '@/api/convpb';
import {clientOptions} from '@/api/grpcweb';
import {pullResource, setValue, trackAction} from '@/api/resource';
import {AccessApiPromiseClient} from '@sc-bos/ui-gen/proto/access_grpc_web_pb';
import {GetLastAccessAttemptRequest, PullAccessAttemptsRequest} from '@sc-bos/ui-gen/proto/access_pb';

/**
 * @param {PullAccessAttemptsRequest.AsObject} request
 * @param {ResourceValue<AccessAttempt.AsObject, PullAccessAttemptsResponse>} resource
 */
export function pullAccessAttempts(request, resource) {
  pullResource('Access.pullAccessAttempts', resource, endpoint => {
    const api = apiClient(endpoint);
    const stream = api.pullAccessAttempts(pullAccessAttemptsRequestFromObject(request));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getAccessAttempt().toObject());
      }
    });
    return stream;
  });
}

/**
 * @param {GetLastAccessAttemptRequest.AsObject} request
 * @param {ActionTracker<AccessAttempt.AsObject>} [tracker]
 * @return {Promise<AccessAttempt.AsObject>}
 */
export function getLastAccessAttempt(request, tracker) {
  return trackAction('Access.getLastAccessAttempt', tracker ?? {}, endpoint => {
    const api = apiClient(endpoint);
    return api.getLastAccessAttempt(getLastAccessAttemptRequestFromObject(request));
  });
}

/**
 * @param {string} endpoint
 * @return {AccessApiPromiseClient}
 */
function apiClient(endpoint) {
  return new AccessApiPromiseClient(endpoint, null, clientOptions());
}

/**
 * @param {PullAccessAttemptsRequest.AsObject} obj
 * @return {PullAccessAttemptsRequest|undefined}
 */
function pullAccessAttemptsRequestFromObject(obj) {
  if (!obj) return undefined;

  const dst = new PullAccessAttemptsRequest();
  setProperties(dst, obj, 'name', 'updatesOnly');
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  return dst;
}

/**
 * @param {GetLastAccessAttemptRequest.AsObject} obj
 * @return {undefined|GetLastAccessAttemptRequest}
 */
function getLastAccessAttemptRequestFromObject(obj) {
  if (!obj) return undefined;

  const dst = new GetLastAccessAttemptRequest();
  setProperties(dst, obj, 'name');
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  return dst;
}
