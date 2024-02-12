import {fieldMaskFromObject, setProperties} from '@/api/convpb';
import {clientOptions} from '@/api/grpcweb';
import {trackAction} from '@/api/resource';
import {MetadataApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/metadata_grpc_web_pb';
import {GetMetadataRequest} from '@smart-core-os/sc-api-grpc-web/traits/metadata_pb';


/**
 * @param {Partial<GetMetadataRequest.AsObject>} request
 * @param {ActionTracker<Metadata.AsObject>} [tracker]
 * @return {Promise<Metadata.AsObject>}
 */
export function getMetadata(request, tracker) {
  return trackAction('Metadata.getMetadata', tracker ?? {}, endpoint => {
    const api = apiClient(endpoint);
    return api.getMetadata(getMetadataRequestFromObject(request));
  });
}

/**
 * @param {string} endpoint
 * @return {MetadataApiPromiseClient}
 */
function apiClient(endpoint) {
  return new MetadataApiPromiseClient(endpoint, null, clientOptions());
}

/**
 * @param {Partial<GetMetadataRequest.AsObject>} obj
 * @return {GetMetadataRequest}
 */
function getMetadataRequestFromObject(obj) {
  if (!obj) return undefined;

  const req = new GetMetadataRequest();
  setProperties(req, obj, 'name');
  req.setReadMask(fieldMaskFromObject(obj.readMask));
  return req;
}
