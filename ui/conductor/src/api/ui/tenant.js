import {clientOptions, simpleFromObject, timestampsFromObject} from '@/api/grpcweb.js';
import {trackAction} from '@/api/resource.js';
import {TenantApiPromiseClient} from '@bsp-ew/ui-gen/src/tenants_grpc_web_pb.js';
import {
  CreateSecretRequest,
  DeleteSecretRequest,
  GetTenantRequest,
  ListSecretsRequest,
  ListTenantsRequest,
  Secret,
  Tenant
} from '@bsp-ew/ui-gen/src/tenants_pb.js';

/**
 * @param {any} request
 * @param {ActionTracker<Tenant.AsObject>} tracker
 * @return {Promise<Tenant.AsObject>}
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
  const proto = new Secret();
  simpleFromObject(proto, obj, 'id', 'secret', 'etag', 'note', 'secretHash');
  timestampsFromObject(proto, obj, 'expireTime', 'firstUseTime', 'lastUseTime');
  if (obj.tenant) {
    proto.setTenant(tenantFromObject(obj.tenant));
  }
  return proto;
}

/**
 * @param {Tenant.AsObject} obj
 * @return {Tenant}
 */
function tenantFromObject(obj) {
  const proto = new Tenant();
  simpleFromObject(proto, obj, 'id', 'title', 'etag', 'zoneNamesList');
  timestampsFromObject(proto, obj, 'createTime');
  return proto;
}
