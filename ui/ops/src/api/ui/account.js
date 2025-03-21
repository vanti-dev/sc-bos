import {convertProperties, fieldMaskFromObject, setProperties, timestampFromObject} from '@/api/convpb.js';
import {clientOptions} from '@/api/grpcweb.js';
import {trackAction} from '@/api/resource.js';
import {AccountApiPromiseClient, AccountInfoPromiseClient} from '@vanti-dev/sc-bos-ui-gen/proto/account_grpc_web_pb';
import {
  Account,
  CreateAccountRequest,
  CreateRoleAssignmentRequest,
  CreateRoleRequest,
  CreateServiceCredentialRequest,
  DeleteAccountRequest,
  DeleteRoleRequest,
  DeleteServiceCredentialRequest,
  GetAccountLimitsRequest,
  GetAccountRequest,
  GetPermissionRequest,
  GetRoleAssignmentRequest,
  GetRoleRequest,
  GetServiceCredentialRequest,
  ListAccountsRequest,
  ListPermissionsRequest,
  ListRoleAssignmentsRequest,
  ListRolesRequest,
  ListServiceCredentialsRequest,
  Role,
  RoleAssignment,
  ServiceCredential,
  UpdateAccountPasswordRequest,
  UpdateAccountRequest,
  UpdateRoleRequest
} from '@vanti-dev/sc-bos-ui-gen/proto/account_pb';

/**
 * @param {Partial<GetAccountRequest.AsObject>} request
 * @param {ActionTracker<Account.AsObject>} [tracker]
 * @return {Promise<Account.AsObject>}
 */
export function getAccount(request, tracker = {}) {
  return trackAction('Account.getAccount', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.getAccount(getAccountRequestFromObject(request));
  });
}

/**
 * @param {Partial<ListAccountsRequest.AsObject>} request
 * @param {ActionTracker<ListAccountsResponse.AsObject>} [tracker]
 * @return {Promise<ListAccountsResponse.AsObject>}
 */
export function listAccounts(request, tracker = {}) {
  return trackAction('Account.listAccounts', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.listAccounts(listAccountsRequestFromObject(request));
  });
}

/**
 * @param {Partial<CreateAccountRequest.AsObject>} request
 * @param {ActionTracker<Account.AsObject>} [tracker]
 * @return {Promise<Account.AsObject>}
 */
export function createAccount(request, tracker = {}) {
  return trackAction('Account.createAccount', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.createAccount(createAccountRequestFromObject(request));
  });
}

/**
 * @param {Partial<UpdateAccountRequest.AsObject>} request
 * @param {ActionTracker<Account.AsObject>} [tracker]
 * @return {Promise<Account.AsObject>}
 */
export function updateAccount(request, tracker = {}) {
  return trackAction('Account.updateAccount', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.updateAccount(updateAccountRequestFromObject(request));
  });
}

/**
 * @param {Partial<DeleteAccountRequest.AsObject>} request
 * @param {ActionTracker<DeleteAccountResponse.AsObject>} [tracker]
 * @return {Promise<DeleteAccountResponse.AsObject>}
 */
export function deleteAccount(request, tracker = {}) {
  return trackAction('Account.deleteAccount', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.deleteAccount(deleteAccountRequestFromObject(request));
  });
}

/**
 * @param {Partial<UpdateAccountPasswordRequest.AsObject>} request
 * @param {ActionTracker<UpdateAccountPasswordResponse.AsObject>} [tracker]
 * @return {Promise<UpdateAccountPasswordResponse.AsObject>}
 */
export function updateAccountPassword(request, tracker = {}) {
  return trackAction('Account.updateAccountPassword', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.updateAccountPassword(updateAccountPasswordRequestFromObject(request));
  });
}

/**
 * @param {Partial<GetServiceCredentialRequest.AsObject>} request
 * @param {ActionTracker<ServiceCredential.AsObject>} [tracker]
 * @return {Promise<ServiceCredential.AsObject>}
 */
export function getServiceCredential(request, tracker = {}) {
  return trackAction('Account.getServiceCredential', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.getServiceCredential(getServiceCredentialRequestFromObject(request));
  });
}

/**
 * @param {Partial<ListServiceCredentialsRequest.AsObject>} request
 * @param {ActionTracker<ListServiceCredentialsResponse.AsObject>} [tracker]
 * @return {Promise<ListServiceCredentialsResponse.AsObject>}
 */
export function listServiceCredentials(request, tracker = {}) {
  return trackAction('Account.listServiceCredentials', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.listServiceCredentials(listServiceCredentialsRequestFromObject(request));
  });
}

/**
 * @param {Partial<CreateServiceCredentialRequest.AsObject>} request
 * @param {ActionTracker<ServiceCredential.AsObject>} [tracker]
 * @return {Promise<ServiceCredential.AsObject>}
 */
export function createServiceCredential(request, tracker = {}) {
  return trackAction('Account.createServiceCredential', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.createServiceCredential(createServiceCredentialRequestFromObject(request));
  });
}

/**
 * @param {Partial<DeleteServiceCredentialRequest.AsObject>} request
 * @param {ActionTracker<DeleteServiceCredentialResponse.AsObject>} [tracker]
 * @return {Promise<DeleteServiceCredentialResponse.AsObject>}
 */
export function deleteServiceCredential(request, tracker = {}) {
  return trackAction('Account.deleteServiceCredential', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.deleteServiceCredential(deleteServiceCredentialRequestFromObject(request));
  });
}

/**
 * @param {Partial<GetRoleRequest.AsObject>} request
 * @param {ActionTracker<Role.AsObject>} [tracker]
 * @return {Promise<Role.AsObject>}
 */
export function getRole(request, tracker = {}) {
  return trackAction('Account.getRole', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.getRole(getRoleRequestFromObject(request));
  });
}

/**
 * @param {Partial<ListRolesRequest.AsObject>} request
 * @param {ActionTracker<ListRolesResponse.AsObject>} [tracker]
 * @return {Promise<ListRolesResponse.AsObject>}
 */
export function listRoles(request, tracker = {}) {
  return trackAction('Account.listRoles', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.listRoles(listRolesRequestFromObject(request));
  });
}

/**
 * @param {Partial<CreateRoleRequest.AsObject>} request
 * @param {ActionTracker<Role.AsObject>} [tracker]
 * @return {Promise<Role.AsObject>}
 */
export function createRole(request, tracker = {}) {
  return trackAction('Account.createRole', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.createRole(createRoleRequestFromObject(request));
  });
}

/**
 * @param {Partial<UpdateRoleRequest.AsObject>} request
 * @param {ActionTracker<Role.AsObject>} [tracker]
 * @return {Promise<Role.AsObject>}
 */
export function updateRole(request, tracker = {}) {
  return trackAction('Account.updateRole', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.updateRole(updateRoleRequestFromObject(request));
  });
}

/**
 * @param {Partial<DeleteRoleRequest.AsObject>} request
 * @param {ActionTracker<DeleteRoleResponse.AsObject>} [tracker]
 * @return {Promise<DeleteRoleResponse.AsObject>}
 */
export function deleteRole(request, tracker = {}) {
  return trackAction('Account.deleteRole', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.deleteRole(deleteRoleRequestFromObject(request));
  });
}

/**
 * @param {Partial<GetRoleAssignmentRequest.AsObject>} request
 * @param {ActionTracker<RoleAssignment.AsObject>} [tracker]
 * @return {Promise<RoleAssignment.AsObject>}
 */
export function getRoleAssignment(request, tracker = {}) {
  return trackAction('Account.getRoleAssignment', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.getRoleAssignment(getRoleAssignmentRequestFromObject(request));
  });
}

/**
 * @param {Partial<ListRoleAssignmentsRequest.AsObject>} request
 * @param {ActionTracker<ListRoleAssignmentsResponse.AsObject>} [tracker]
 * @return {Promise<ListRoleAssignmentsResponse.AsObject>}
 */
export function listRoleAssignments(request, tracker = {}) {
  return trackAction('Account.listRoleAssignments', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.listRoleAssignments(listRoleAssignmentsRequestFromObject(request));
  });
}

/**
 * @param {Partial<CreateRoleAssignmentRequest.AsObject>} request
 * @param {ActionTracker<RoleAssignment.AsObject>} [tracker]
 * @return {Promise<RoleAssignment.AsObject>}
 */
export function createRoleAssignment(request, tracker = {}) {
  return trackAction('Account.createRoleAssignment', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.createRoleAssignment(createRoleAssignmentRequestFromObject(request));
  });
}

/**
 * @param {Partial<DeleteRoleAssignmentRequest.AsObject>} request
 * @param {ActionTracker<DeleteRoleAssignmentResponse.AsObject>} [tracker]
 * @return {Promise<DeleteRoleAssignmentResponse.AsObject>}
 */
export function deleteRoleAssignment(request, tracker = {}) {
  return trackAction('Account.deleteRoleAssignment', tracker, endpoint => {
    const api = apiClient(endpoint);
    return api.deleteRoleAssignment(deleteRoleAssignmentRequestFromObject(request));
  });
}

/**
 * @param {Partial<GetPermissionRequest.AsObject>} request
 * @param {ActionTracker<Permission.AsObject>} [tracker]
 * @return {Promise<Permission.AsObject>}
 */
export function getPermission(request, tracker = {}) {
  return trackAction('Account.getPermission', tracker, endpoint => {
    const api = infoClient(endpoint);
    return api.getPermission(getPermissionRequestFromObject(request));
  });
}

/**
 * @param {Partial<ListPermissionsRequest.AsObject>} request
 * @param {ActionTracker<ListPermissionsResponse.AsObject>} [tracker]
 * @return {Promise<ListPermissionsResponse.AsObject>}
 */
export function listPermissions(request, tracker = {}) {
  return trackAction('Account.listPermissions', tracker, endpoint => {
    const api = infoClient(endpoint);
    return api.listPermissions(listPermissionsRequestFromObject(request));
  });
}

/**
 * @param {Partial<GetAccountLimitsRequest.AsObject>} request
 * @param {ActionTracker<AccountLimits.AsObject>} [tracker]
 * @return {Promise<AccountLimits.AsObject>}
 */
export function getAccountLimits(request, tracker = {}) {
  return trackAction('Account.getAccountLimits', tracker, endpoint => {
    const api = infoClient(endpoint);
    return api.getAccountLimits(getAccountLimitsRequestFromObject(request));
  });
}

/**
 * @param {string} endpoint
 * @return {AccountApiPromiseClient}
 */
function apiClient(endpoint) {
  return new AccountApiPromiseClient(endpoint, null, clientOptions());
}

/**
 * @param {string} endpoint
 * @return {AccountInfoPromiseClient}
 */
function infoClient(endpoint) {
  return new AccountInfoPromiseClient(endpoint, null, clientOptions());
}

/**
 * @param {Partial<GetAccountRequest.AsObject>} obj
 * @return {undefined|GetAccountRequest}
 */
function getAccountRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new GetAccountRequest();
  setProperties(dst, obj, 'name', 'id');
  return dst;
}

/**
 * @param {Partial<ListAccountsRequest.AsObject>} obj
 * @return {undefined|ListAccountsRequest}
 */
function listAccountsRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new ListAccountsRequest();
  setProperties(dst, obj, 'name', 'pageToken', 'pageSize');
  return dst;
}

/**
 * @param {Partial<CreateAccountRequest.AsObject>} obj
 * @return {undefined|CreateAccountRequest}
 */
function createAccountRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new CreateAccountRequest();
  setProperties(dst, obj, 'name', 'password');
  dst.setAccount(accountFromObject(obj.account));
  return dst;
}

/**
 * @param {Partial<UpdateAccountRequest.AsObject>} obj
 * @return {undefined|UpdateAccountRequest}
 */
function updateAccountRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new UpdateAccountRequest();
  setProperties(dst, obj, 'name');
  dst.setUpdateMask(fieldMaskFromObject(obj.updateMask));
  dst.setAccount(accountFromObject(obj.account));
  return dst;
}

/**
 * @param {Partial<DeleteAccountRequest.AsObject>} obj
 * @return {undefined|DeleteAccountRequest}
 */
function deleteAccountRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new DeleteAccountRequest();
  setProperties(dst, obj, 'name', 'id', 'allowMissing');
  return dst;
}

/**
 * @param {Partial<UpdateAccountPasswordRequest.AsObject>} obj
 * @return {undefined|UpdateAccountPasswordRequest}
 */
function updateAccountPasswordRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new UpdateAccountPasswordRequest();
  setProperties(dst, obj, 'name', 'id', 'oldPassword', 'newPassword');
  return dst;
}

/**
 * @param {Partial<Account.AsObject>} obj
 * @return {undefined|Account}
 */
function accountFromObject(obj) {
  if (!obj) return undefined;
  const dst = new Account();
  setProperties(dst, obj, 'id', 'type', 'displayName', 'description', 'username');
  convertProperties(dst, obj, timestampFromObject, 'createTime');
  return dst;
}

/**
 * @param {Partial<GetServiceCredentialRequest.AsObject>} obj
 * @return {undefined|GetServiceCredentialRequest}
 */
function getServiceCredentialRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new GetServiceCredentialRequest();
  setProperties(dst, obj, 'name', 'id');
  return dst;
}

/**
 * @param {Partial<ListServiceCredentialsRequest.AsObject>} obj
 * @return {undefined|ListServiceCredentialsRequest}
 */
function listServiceCredentialsRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new ListServiceCredentialsRequest();
  setProperties(dst, obj, 'name', 'accountId');
  return dst;
}

/**
 * @param {Partial<CreateServiceCredentialRequest.AsObject>} obj
 * @return {undefined|CreateServiceCredentialRequest}
 */
function createServiceCredentialRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new CreateServiceCredentialRequest();
  setProperties(dst, obj, 'name');
  dst.setServiceCredential(serviceCredentialFromObject(obj.serviceCredential));
  return dst;
}

/**
 * @param {Partial<DeleteServiceCredentialRequest.AsObject>} obj
 * @return {undefined|DeleteServiceCredentialRequest}
 */
function deleteServiceCredentialRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new DeleteServiceCredentialRequest();
  setProperties(dst, obj, 'name', 'id', 'allowMissing');
  return dst;
}

/**
 * @param {Partial<ServiceCredential.AsObject>} obj
 * @return {undefined|ServiceCredential}
 */
function serviceCredentialFromObject(obj) {
  if (!obj) return undefined;
  const dst = new ServiceCredential();
  setProperties(dst, obj, 'id', 'accountId', 'displayName', 'description', 'secret');
  convertProperties(dst, obj, timestampFromObject, 'createTime', 'expireTime');
  return dst;
}

/**
 * @param {Partial<GetRoleRequest.AsObject>} obj
 * @return {undefined|GetRoleRequest}
 */
function getRoleRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new GetRoleRequest();
  setProperties(dst, obj, 'name', 'id');
  return dst;
}

/**
 * @param {Partial<ListRolesRequest.AsObject>} obj
 * @return {undefined|ListRolesRequest}
 */
function listRolesRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new ListRolesRequest();
  setProperties(dst, obj, 'name', 'pageToken', 'pageSize');
  return dst;
}

/**
 * @param {Partial<CreateRoleRequest.AsObject>} obj
 * @return {undefined|CreateRoleRequest}
 */
function createRoleRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new CreateRoleRequest();
  setProperties(dst, obj, 'name');
  dst.setRole(roleFromObject(obj.role));
  return dst;
}

/**
 * @param {Partial<UpdateRoleRequest.AsObject>} obj
 * @return {undefined|UpdateRoleRequest}
 */
function updateRoleRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new UpdateRoleRequest();
  setProperties(dst, obj, 'name');
  dst.setUpdateMask(fieldMaskFromObject(obj.updateMask));
  dst.setRole(roleFromObject(obj.role));
  return dst;
}

/**
 * @param {Partial<DeleteRoleRequest.AsObject>} obj
 * @return {undefined|DeleteRoleRequest}
 */
function deleteRoleRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new DeleteRoleRequest();
  setProperties(dst, obj, 'name', 'id', 'allowMissing');
  return dst;
}

/**
 * @param {Partial<Role.AsObject>} obj
 * @return {undefined|Role}
 */
function roleFromObject(obj) {
  if (!obj) return undefined;
  const dst = new Role();
  setProperties(dst, obj, 'id', 'displayName', 'description');
  dst.setPermissionIdsList(obj.permissionIdsList ?? [])
  return dst;
}

/**
 * @param {Partial<GetRoleAssignmentRequest.AsObject>} obj
 * @return {undefined|GetRoleAssignmentRequest}
 */
function getRoleAssignmentRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new GetRoleAssignmentRequest();
  setProperties(dst, obj, 'name', 'id');
  return dst;
}

/**
 * @param {Partial<ListRoleAssignmentsRequest.AsObject>} obj
 * @return {undefined|ListRoleAssignmentsRequest}
 */
function listRoleAssignmentsRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new ListRoleAssignmentsRequest();
  setProperties(dst, obj, 'name', 'pageToken', 'pageSize', 'filter');
  return dst;
}

/**
 * @param {Partial<CreateRoleAssignmentRequest.AsObject>} obj
 * @return {undefined|CreateRoleAssignmentRequest}
 */
function createRoleAssignmentRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new CreateRoleAssignmentRequest();
  setProperties(dst, obj, 'name');
  dst.setRoleAssignment(roleAssignmentFromObject(obj.roleAssignment));
  return dst;
}

/**
 * @param {Partial<DeleteRoleRequest.AsObject>} obj
 * @return {undefined|DeleteRoleRequest}
 */
function deleteRoleAssignmentRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new DeleteRoleRequest();
  setProperties(dst, obj, 'name', 'id', 'allowMissing');
  return dst;
}

/**
 * @param {Partial<RoleAssignment.AsObject>} obj
 * @return {undefined|RoleAssignment}
 */
function roleAssignmentFromObject(obj) {
  if (!obj) return undefined;
  const dst = new RoleAssignment();
  setProperties(dst, obj, 'id', 'accountId', 'roleId');
  dst.setScope(roleAssignmentScopeFromObject(obj.scope));
  return dst;
}

/**
 * @param {Partial<RoleAssignment.Scope.AsObject>} obj
 * @return {undefined|RoleAssignment.Scope}
 */
function roleAssignmentScopeFromObject(obj) {
  if (!obj) return undefined;
  const dst = new RoleAssignment.Scope();
  setProperties(dst, obj, 'resourceType', 'resource');
  return dst;
}

/**
 * @param {Partial<GetPermissionRequest.AsObject>} obj
 * @return {undefined|GetPermissionRequest}
 */
function getPermissionRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new GetPermissionRequest();
  setProperties(dst, obj, 'name', 'id');
  return dst;
}

/**
 * @param {Partial<ListPermissionsRequest.AsObject>} obj
 * @return {undefined|ListPermissionsRequest}
 */
function listPermissionsRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new ListPermissionsRequest();
  setProperties(dst, obj, 'name', 'pageToken', 'pageSize');
  return dst;
}

/**
 * @param {Partial<GetAccountLimitsRequest.AsObject>} obj
 * @return {undefined|GetAccountLimitsRequest}
 */
function getAccountLimitsRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new GetAccountLimitsRequest();
  setProperties(dst, obj, 'name');
  return dst;
}

/**
 * A map from RoleAssignment.ResourceType to the enum name, the inverse of RoleAssignment.ResourceType.
 *
 * @type {Record<number, keyof RoleAssignment.ResourceType>}
 */
export const ResourceTypeById =
    Object.entries(RoleAssignment.ResourceType).reduce((acc, [k, v]) => {
      acc[v] = k;
      return acc;
    }, {});
