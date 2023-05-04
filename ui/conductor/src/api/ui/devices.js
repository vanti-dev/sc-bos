import {setProperties} from '@/api/convpb.js';
import {clientOptions} from '@/api/grpcweb.js';
import {trackAction} from '@/api/resource.js';
import {DevicesApiPromiseClient} from '@sc-bos/ui-gen/proto/devices_grpc_web_pb';
import {Device, DevicesMetadata, GetDevicesMetadataRequest, ListDevicesRequest} from '@sc-bos/ui-gen/proto/devices_pb';

/**
 * @param {ListDevicesRequest.AsObject} request
 * @param {ActionTracker<ListDevicesResponse.AsObject>} tracker
 * @return {Promise<ListDevicesResponse.AsObject>}
 */
export function listDevices(request, tracker) {
  return trackAction('Devices.listDevices', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.listDevices(listDevicesRequestFromObject(request));
  });
}

/**
 * @param {string} endpoint
 * @return {DevicesApiPromiseClient}
 */
function client(endpoint) {
  return new DevicesApiPromiseClient(endpoint, null, clientOptions());
}

/**
 * @param {ListDevicesRequest.AsObject} obj
 * @return {undefined|ListDevicesRequest}
 */
function listDevicesRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new ListDevicesRequest();
  setProperties(dst, obj, 'pageToken', 'pageSize');
  dst.setQuery(deviceQueryFromObject(obj.query));
  return dst;
}

/**
 * @param {Device.Query.AsObject} obj
 * @return {undefined|Device.Query}
 */
function deviceQueryFromObject(obj) {
  if (!obj) return undefined;
  const dst = new Device.Query();
  for (const item of (obj.conditionsList ?? [])) {
    const dstItem = new Device.Query.Condition();
    setProperties(dstItem, item, 'field',
        'stringEqual', 'stringEqualFold',
        'stringContains', 'stringContainsFold'
    );
    dst.addConditions(dstItem);
  }
  return dst;
}

/**
 *
 * @param {string} endpoint
 * @param {ActionTracker<GetDevicesMetadataResponse.AsObject>} tracker
 * @return {Promise<GetDevicesMetadataResponse.AsObject>}
 */
export function getDevicesMetadata(endpoint, tracker) {
  return trackAction('Devices.getDevicesMetadata', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    const request = new GetDevicesMetadataRequest();
    const includes = new DevicesMetadata.Include();
    includes.addFields('metadata.membership.subsystem');
    request.setIncludes(includes);
    return api.getDevicesMetadata(request);
  });
};
