import {trackAction} from '@/api/resource';
import {ServicesApiPromiseClient} from '@sc-bos/ui-gen/proto/services_grpc_web_pb';
import {clientOptions} from '@/api/grpcweb';
import {GetMetadataRequest} from '@smart-core-os/sc-api-grpc-web/traits/metadata_pb';
import {fieldMaskFromObject, setProperties} from '@/api/convpb';


/**
 * @param {GetMetadataRequest.AsObject} request
 * @param {ActionTracker<ServiceMetadata.AsObject>} tracker
 * @return {Promise<ServiceMetadata.AsObject>}
 */
export function getServiceMetadata(request, tracker) {
  const name = String(request.name);
  if (!name) throw new Error('request.name must be specified');
  return trackAction('Services.GetServiceMetadata', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.getServiceMetadata(createGetMetadataRequestFromObject(request));
  });
}


/**
 * @param {string} endpoint
 * @return {ServicesApiPromiseClient}
 */
function client(endpoint) {
  return new ServicesApiPromiseClient(endpoint, null, clientOptions());
}

/**
 * @param {GetMetadataRequest.AsObject} obj
 * @return {GetMetadataRequest}
 */
function createGetMetadataRequestFromObject(obj) {
  if (!obj) return undefined;

  const req = new GetMetadataRequest();
  setProperties(req, obj, 'name');
  req.setReadMask(fieldMaskFromObject(obj.readMask));
  return req;
}

export const ServiceNames = {
  Automations: 'automations',
  Drivers: 'drivers',
  Systems: 'systems'
};
