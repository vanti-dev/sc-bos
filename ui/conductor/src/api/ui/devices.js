import {fieldMaskFromObject, setProperties} from '@/api/convpb.js';
import {clientOptions} from '@/api/grpcweb.js';
import {pullResource, setValue} from '@/api/resource';
import {trackAction} from '@/api/resource.js';
import {DevicesApiPromiseClient} from '@sc-bos/ui-gen/proto/devices_grpc_web_pb';
import {
  Device,
  DevicesMetadata,
  GetDevicesMetadataRequest,
  ListDevicesRequest,
  PullDevicesMetadataRequest
} from '@sc-bos/ui-gen/proto/devices_pb';

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
 * @param {GetDevicesMetadataRequest.AsObject} request
 * @param {ActionTracker<GetDevicesMetadataRequest.AsObject>} tracker
 * @return {Promise<DevicesMetadata.AsObject>}
 */
export function getDevicesMetadata(request, tracker) {
  return trackAction('Devices.getDevicesMetadata', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.getDevicesMetadata(getDevicesMetadataRequestFromObject(request));
  });
}

/**
 * @param {PullDevicesMetadataRequest.AsObject} request
 * @param {ResourceValue<DevicesMetadata.AsObject, DevicesMetadata>} resource
 */
export function pullDevicesMetadata(request, resource) {
  pullResource('Devices.pullDevicesMetadata', resource, endpoint => {
    const api = client(endpoint);
    const stream = api.pullDevicesMetadata(pullDevicesMetadataRequestFromObject(request));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getDevicesMetadata().toObject());
      }
    });
    return stream;
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
 * @param {DevicesMetadata.Include.AsObject} obj
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

/**
 *
 * @param {PullDevicesMetadataRequest.AsObject} obj
 * @return {undefined|PullDevicesMetadataRequest}
 */
function pullDevicesMetadataRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new PullDevicesMetadataRequest();
  setProperties(dst, obj, 'updatesOnly');
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  dst.setIncludes(devicesMetadataIncludeFromObject(obj.includes));
  return dst;
}
