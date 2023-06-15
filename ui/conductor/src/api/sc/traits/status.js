import {fieldMaskFromObject, setProperties} from '@/api/convpb';
import {clientOptions} from '@/api/grpcweb';
import {pullResource, setValue, trackAction} from '@/api/resource';
import {StatusApiPromiseClient} from '@sc-bos/ui-gen/proto/status_grpc_web_pb';
import {GetCurrentStatusRequest, PullCurrentStatusRequest} from '@sc-bos/ui-gen/proto/status_pb';


/**
 * @param {PullCurrentStatusRequest.AsObject} request
 * @param {ResourceValue<StatusLog.AsObject, PullCurrentStatusResponse>} resource
 */
export function pullCurrentStatus(request, resource) {
  pullResource('Status.pullCurrentStatus', resource, endpoint => {
    const api = new StatusApiPromiseClient(endpoint, null, clientOptions());
    const stream = api.pullCurrentStatus(pullCurrentStatusRequestFromObject(request));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getCurrentStatus().toObject());
      }
    });
    return stream;
  });
}

/**
 * @param {GetCurrentStatusRequest.AsObject} request
 * @param {ActionTracker<StatusLog.AsObject>} tracker
 * @return {Promise<StatusLog.AsObject>}
 */
export function getCurrentStatus(request, tracker) {
  return trackAction('Status.getCurrentStatus', tracker ?? {}, endpoint => {
    const api = new StatusApiPromiseClient(endpoint, null, clientOptions());
    return api.getCurrentStatus(getCurrentStatusRequestFromObject(request));
  });
}

/**
 * @param {PullCurrentStatusRequest.AsObject} obj
 * @return {PullCurrentStatusRequest|undefined}
 */
function pullCurrentStatusRequestFromObject(obj) {
  if (!obj) return undefined;

  const dst = new PullCurrentStatusRequest();
  setProperties(dst, obj, 'name', 'updatesOnly');
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  return dst;
}

/**
 * @param {GetCurrentStatusRequest.AsObject} obj
 * @return {undefined|GetCurrentStatusRequest}
 */
function getCurrentStatusRequestFromObject(obj) {
  if (!obj) return undefined;

  const dst = new GetCurrentStatusRequest();
  setProperties(dst, obj, 'name');
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  return dst;
}
