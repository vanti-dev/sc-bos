import {fieldMaskFromObject, setProperties} from '@/api/convpb.js';
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
 *
 * @param {Promise<GetDevicesMetadataRequest.AsObject>} request
 * @param {ActionTracker<GetDevicesMetadataResponse.AsObject>} tracker
 * @return {Promise<GetDevicesMetadataResponse.AsObject>}
 */
export function getDevicesMetadata(request, tracker) {
  return trackAction('Devices.getDevicesMetadata', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.getDevicesMetadata(getDevicesMetadataRequestFromObject(request));
  });
};

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
 * @param {GetDevicesMetadataRequest.AsObject} obj
 * @return {undefined|DevicesMetadata.Include}
 */
function devicesMetadataIncludeFromObject(obj) {
  if (!obj) return undefined;
  const dst = new DevicesMetadata.Include();
  dst.setFieldsList(obj.fieldsList);
  return dst;
}

/**
 *
 * @param {GetDevicesMetadataRequest.AsObject} obj
 * @return {undefined|GetDevicesMetadataRequest}
 */
function getDevicesMetadataRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new GetDevicesMetadataRequest();
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  dst.setIncludes(devicesMetadataIncludeFromObject(obj.includes));
  return dst;
}
