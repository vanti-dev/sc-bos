import {convertProperties, setProperties, timestampFromObject, timestampToDate} from '@/api/convpb.js';
import {clientOptions} from '@/api/grpcweb.js';
import {trackAction} from '@/api/resource.js';
import {TenantApiPromiseClient} from '@sc-bos/ui-gen/proto/tenants_grpc_web_pb';
import {
  AddTenantZonesRequest,
  CreateSecretRequest,
  CreateTenantRequest,
  DeleteSecretRequest,
  DeleteTenantRequest,
  GetTenantRequest,
  ListSecretsRequest,
  ListTenantsRequest,
  RemoveTenantZonesRequest,
  Secret,
  Tenant
} from '@sc-bos/ui-gen/proto/tenants_pb';

/**
 * @param {ListTenantsRequest.AsObject} request
 * @param {ActionTracker<ListTenantsResponse.AsObject>} [tracker]
 * @return {Promise<ListTenantsResponse.AsObject>}
 */
export function listTenants(request, tracker) {
  return trackAction('Tenant.listTenants', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.listTenants(new ListTenantsRequest());
  });
}

/**
 *
 * @param {CreateTenantRequest.AsObject} request
 * @param {ActionTracker<CreateTenantRequest.AsObject>} [tracker]
 * @return {Promise<Tenant.AsObject>}
 */
export function createTenant(request, tracker) {
  return trackAction('Tenant.createTenant', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.createTenant(createTenantRequestFromObject(request));
  });
}

/**
 *
 * @param {DeleteTenantRequest.AsObject} obj
 * @param {ActionTracker<DeleteTenantRequest.AsObject>} [tracker]
 * @return {Promise<DeleteTenantResponse>}
 */
export function deleteTenant(obj, tracker) {
  return trackAction('Tenant.deleteTenant', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.deleteTenant(deleteTenantRequestFromObject(obj));
  });
}

/**
 *
 * @param {GetTenantRequest.AsObject} request
 * @param {ActionTracker<GetTenantRequest.AsObject>} [tracker]
 * @return {Promise<Tenant.AsObject>}
 */
export function getTenant(request, tracker) {
  const id = String(request.id);
  if (!id) throw new Error('request.id must be specified');
  return trackAction('Tenant.getTenant', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.getTenant(new GetTenantRequest().setId(id));
  });
}

/**
 *
 * @param {ListSecretsRequest.AsObject} request
 * @param {ActionTracker<ListSecretsRequest.AsObject>} [tracker]
 * @return {Promise<ListSecretsResponse.AsObject>}
 */
export function listSecrets(request, tracker) {
  const tenantId = String(request.tenantId);
  if (!tenantId) throw new Error('request.tenantId must be specified');
  return trackAction('Tenant.listSecrets', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.listSecrets(new ListSecretsRequest().setFilter(`tenant.id=${tenantId}`));
  });
}

/**
 *
 * @param {CreateSecretRequest.AsObject} request
 * @param {ActionTracker<CreateSecretRequest.AsObject>} [tracker]
 * @return {Promise<Secret.AsObject>}
 */
export function createSecret(request, tracker) {
  const secret = request.secret;
  if (!secret) throw new Error('request.secret must be specified');
  return trackAction('Tenant.createSecret', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.createSecret(new CreateSecretRequest().setSecret(secretFromObject(secret)));
  });
}

/**
 *
 * @param {DeleteSecretRequest.AsObject} request
 * @param {ActionTracker<DeleteSecretRequest.AsObject>} [tracker]
 * @return {Promise<DeleteSecretResponse.AsObject>}
 */
export function deleteSecret(request, tracker) {
  const secretId = request.id;
  if (!secretId) throw new Error('request.id must be specified');
  return trackAction('Tenant.deleteSecret', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.deleteSecret(new DeleteSecretRequest().setId(secretId));
  });
}

/**
 * @param {AddTenantZonesRequest.AsObject} request
 * @param {ActionTracker<AddTenantZonesRequest.AsObject>} [tracker]
 * @return {Promise<Tenant.AsObject>}
 */
export function addTenantZones(request, tracker) {
  return trackAction('Tenant.addTenantZones', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.addTenantZones(addTenantZonesRequestFromObject(request));
  });
}

/**
 * @param {RemoveTenantZonesRequest.AsObject} request
 * @param {ActionTracker<RemoveTenantZonesRequest.AsObject>} [tracker]
 * @return {Promise<Tenant.AsObject>}
 */
export function removeTenantZones(request, tracker) {
  return trackAction('Tenant.removeTenantZones', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.removeTenantZones(removeTenantZonesRequestFromObject(request));
  });
}

/**
 * @param {string} endpoint
 * @return {TenantApiPromiseClient}
 */
function client(endpoint) {
  return new TenantApiPromiseClient(endpoint, null, clientOptions());
}

/**
 * @param {CreateTenantRequest.AsObject} obj
 * @return {CreateTenantRequest}
 */
function createTenantRequestFromObject(obj) {
  if (!obj) return undefined;

  const req = new CreateTenantRequest();
  req.setTenant(tenantFromObject(obj.tenant));
  return req;
}

/**
 *
 * @param {DeleteTenantRequest.AsObject} obj
 * @return {DeleteTenantRequest}
 */
function deleteTenantRequestFromObject(obj) {
  if (!obj) return undefined;

  const req = new DeleteTenantRequest();
  setProperties(req, obj, 'id');
  return req;
}

/**
 *
 * @param {Secret.AsObject} obj
 * @return {Secret}
 */
function secretFromObject(obj) {
  if (!obj) return undefined;
  const dst = new Secret();
  setProperties(dst, obj, 'id', 'secret', 'etag', 'note', 'secretHash');
  convertProperties(dst, obj, timestampFromObject, 'expireTime', 'firstUseTime', 'lastUseTime');
  dst.setTenant(tenantFromObject(obj.tenant));
  return dst;
}

/**
 * @param {Secret|Secret.AsObject|null} s
 * @return {Secret.AsObject&{createTime?: Date, expireTime?:Date, lastUseTime?:Date, firstUseTime?:Date} | null}
 */
export function secretToObject(s) {
  if (!s) return null;

  const res = {...s};
  for (const prop of ['createTime', 'expireTime', 'firstUseTime', 'lastUseTime']) {
    if (s[prop]) {
      res[prop] = timestampToDate(s[prop]);
    }
  }
  return res;
}

/**
 * @param {Tenant.AsObject} obj
 * @return {Tenant}
 */
function tenantFromObject(obj) {
  if (!obj) return undefined;
  const dst = new Tenant();
  setProperties(dst, obj, 'id', 'title', 'etag', 'zoneNamesList');
  dst.setCreateTime(timestampFromObject(obj.createTime));
  return dst;
}

/**
 * @param {AddTenantZonesRequest.AsObject} obj
 * @return {undefined|AddTenantZonesRequest}
 */
function addTenantZonesRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new AddTenantZonesRequest();
  setProperties(dst, obj, 'tenantId', 'addZoneNamesList');
  return dst;
}

/**
 * @param {RemoveTenantZonesRequest.AsObject} obj
 * @return {undefined|RemoveTenantZonesRequest}
 */
function removeTenantZonesRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new RemoveTenantZonesRequest();
  setProperties(dst, obj, 'tenantId', 'removeZoneNamesList');
  return dst;
}
