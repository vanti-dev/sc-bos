import {pullResource, setCollection, trackAction} from '@/api/resource';
import {ServicesApiPromiseClient} from '@sc-bos/ui-gen/proto/services_grpc_web_pb';
import {clientOptions} from '@/api/grpcweb';
import {GetMetadataRequest} from '@smart-core-os/sc-api-grpc-web/traits/metadata_pb';
import {fieldMaskFromObject, setProperties} from '@/api/convpb';
import {
  ConfigureServiceRequest,
  ListServicesRequest, PullServiceRequest,
  StartServiceRequest,
  StopServiceRequest
} from '@sc-bos/ui-gen/proto/services_pb';


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
 * @param {ListServicesRequest.AsObject} request
 * @param {ActionTracker<ListServicesRequest.AsObject>} tracker
 * @return {Promise<ListServicesResponse.AsObject>}
 */
export function listServices(request, tracker) {
  const name = String(request.name);
  if (!name) throw new Error('request.name must be specified');
  return trackAction('Services.ListServices', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.listServices(createListServicesRequestFromObject(request));
  });
}

/**
 *
 * @param {PullServiceRequest.AsObject} request
 * @param {ResourceCollection<Service.AsObject, Service>} resource
 */
export function pullServices(request, resource) {
  pullResource('Services.PullServices', resource, endpoint => {
    const api = client(endpoint);
    const stream = api.pullServices(pullServicesRequestFromObject(request));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setCollection(resource, change, v => v.id);
      }
    });
    return stream;
  });
}

/**
 * @param {ConfigureServiceRequest.AsObject} request
 * @param {ActionTracker<Service.AsObject>} tracker
 * @return {Promise<Service.AsObject>}
 */
export function configureService(request, tracker) {
  if (!(request.name && request.id)) throw new Error('request.name and request.id must be specified');
  return trackAction('Services.ConfigureService', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.configureService(createConfigureServiceRequestFromObject(request));
  });
}

/**
 * @param {StartServiceRequest.AsObject} request
 * @param {ActionTracker<Service.AsObject>} tracker
 * @return {Promise<Service.AsObject>}
 */
export function startService(request, tracker) {
  if (!(request.name && request.id)) throw new Error('request.name and request.id must be specified');
  return trackAction('Services.StartService', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.startService(createStartServiceRequestFromObject(request));
  });
}
/**
 * @param {StopServiceRequest.AsObject} request
 * @param {ActionTracker<Service.AsObject>} tracker
 * @return {Promise<Service.AsObject>}
 */
export function stopService(request, tracker) {
  if (!(request.name && request.id)) throw new Error('request.name and request.id must be specified');
  return trackAction('Services.StopService', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.stopService(createStopServiceRequestFromObject(request));
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

/**
 * @param {ListServicesRequest.AsObject} obj
 * @return {ListServicesRequest}
 */
function createListServicesRequestFromObject(obj) {
  if (!obj) return undefined;
  const req = new ListServicesRequest();
  setProperties(req, obj, 'name', 'pageToken', 'pageSize');
  return req;
}

/**
 * @param {PullServiceRequest.AsObject} obj
 * @return {PullServiceRequest}
 */
function pullServicesRequestFromObject(obj) {
  if (!obj) return undefined;
  const req = new PullServiceRequest();
  setProperties(req, obj, 'name', 'updatesOnly');
  req.setReadMask(fieldMaskFromObject(obj.readMask));
  return req;
}

/**
 * @param {ConfigureServiceRequest.AsObject} obj
 * @return {ConfigureServiceRequest}
 */
function createConfigureServiceRequestFromObject(obj) {
  if (!obj) return undefined;
  const req = new ConfigureServiceRequest();
  setProperties(req, obj, 'name', 'id');
  req.setConfigRaw(obj.configRaw);
  return req;
}

/**
 * @param {StartServiceRequest.AsObject} obj
 * @return {StartServiceRequest}
 */
function createStartServiceRequestFromObject(obj) {
  if (!obj) return undefined;
  const req = new StartServiceRequest();
  setProperties(req, obj, 'name', 'id', 'allowActive');
  return req;
}

/**
 * @param {StopServiceRequest.AsObject} obj
 * @return {StopServiceRequest}
 */
function createStopServiceRequestFromObject(obj) {
  if (!obj) return undefined;
  const req = new StopServiceRequest();
  setProperties(req, obj, 'name', 'id', 'allowActive');
  return req;
}


export const ServiceNames = {
  Automations: 'automations',
  Drivers: 'drivers',
  Systems: 'systems',
  Zones: 'zones'
};
