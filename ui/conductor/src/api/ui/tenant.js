import {convertProperties, setProperties, timestampFromObject, timestampToDate} from '@/api/convpb.js';
import {clientOptions} from '@/api/grpcweb.js';
import {trackAction} from '@/api/resource.js';
import {TenantApiPromiseClient} from '@sc-bos/ui-gen/proto/tenants_grpc_web_pb';
import {
  CreateSecretRequest,
  DeleteSecretRequest,
  GetTenantRequest,
  ListSecretsRequest,
  ListTenantsRequest,
  Secret,
  Tenant
} from '@sc-bos/ui-gen/proto/tenants_pb';

/**
 * @param {ListTenantsRequest.AsObject} request
 * @param {ActionTracker<ListTenantsResponse.AsObject>} tracker
 * @return {Promise<ListTenantsResponse.AsObject>}
 */
export function listTenants(request, tracker) {
  return trackAction('Tenant.listTenants', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.listTenants(new ListTenantsRequest());
  })
}

export function getTenant(request, tracker) {
  const id = String(request.id);
  if (!id) throw new Error('request.id must be specified');
  return trackAction('Tenant.getTenant', tracker ?? {}, async endpoint => {
    const api = client(endpoint);
    return api.getTenant(new GetTenantRequest().setId(id));
  })
}

export function listSecrets(request, tracker) {
  const tenantId = String(request.tenantId);
  if (!tenantId) throw new Error('request.tenantId must be specified');
  return trackAction('Tenant.listSecrets', tracker ?? {}, async endpoint => {
    const api = client(endpoint);
    return api.listSecrets(new ListSecretsRequest().setFilter(`tenant.id=${tenantId}`));
  })
}

export function createSecret(request, tracker) {
  const secret = request.secret;
  if (!secret) throw new Error('request.secret must be specified');
  return trackAction('Tenant.createSecret', tracker ?? {}, async endpoint => {
    const api = client(endpoint);
    return api.createSecret(new CreateSecretRequest().setSecret(secretFromObject(secret)));
  });
}

export function deleteSecret(request, tracker) {
  const secretId = request.id;
  if (!secretId) throw new Error('request.id must be specified');
  return trackAction('Tenant.deleteSecret', tracker ?? {}, async endpoint => {
    const api = client(endpoint);
    return api.deleteSecret(new DeleteSecretRequest().setId(secretId));
  })
}

/**
 * @param {string} endpoint
 * @return {TenantApiPromiseClient}
 */
function client(endpoint) {
  return new TenantApiPromiseClient(endpoint, null, clientOptions());
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
