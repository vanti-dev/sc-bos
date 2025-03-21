import * as grpcWeb from 'grpc-web';

import * as account_pb from './account_pb'; // proto import: "account.proto"


export class AccountApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getAccount(
    request: account_pb.GetAccountRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.Account) => void
  ): grpcWeb.ClientReadableStream<account_pb.Account>;

  listAccounts(
    request: account_pb.ListAccountsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.ListAccountsResponse) => void
  ): grpcWeb.ClientReadableStream<account_pb.ListAccountsResponse>;

  createAccount(
    request: account_pb.CreateAccountRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.Account) => void
  ): grpcWeb.ClientReadableStream<account_pb.Account>;

  updateAccount(
    request: account_pb.UpdateAccountRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.Account) => void
  ): grpcWeb.ClientReadableStream<account_pb.Account>;

  updateAccountPassword(
    request: account_pb.UpdateAccountPasswordRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.UpdateAccountPasswordResponse) => void
  ): grpcWeb.ClientReadableStream<account_pb.UpdateAccountPasswordResponse>;

  rotateAccountClientSecret(
    request: account_pb.RotateAccountClientSecretRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.RotateAccountClientSecretResponse) => void
  ): grpcWeb.ClientReadableStream<account_pb.RotateAccountClientSecretResponse>;

  deleteAccount(
    request: account_pb.DeleteAccountRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.DeleteAccountResponse) => void
  ): grpcWeb.ClientReadableStream<account_pb.DeleteAccountResponse>;

  getRole(
    request: account_pb.GetRoleRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.Role) => void
  ): grpcWeb.ClientReadableStream<account_pb.Role>;

  listRoles(
    request: account_pb.ListRolesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.ListRolesResponse) => void
  ): grpcWeb.ClientReadableStream<account_pb.ListRolesResponse>;

  createRole(
    request: account_pb.CreateRoleRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.Role) => void
  ): grpcWeb.ClientReadableStream<account_pb.Role>;

  updateRole(
    request: account_pb.UpdateRoleRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.Role) => void
  ): grpcWeb.ClientReadableStream<account_pb.Role>;

  deleteRole(
    request: account_pb.DeleteRoleRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.DeleteRoleResponse) => void
  ): grpcWeb.ClientReadableStream<account_pb.DeleteRoleResponse>;

  getRoleAssignment(
    request: account_pb.GetRoleAssignmentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.RoleAssignment) => void
  ): grpcWeb.ClientReadableStream<account_pb.RoleAssignment>;

  listRoleAssignments(
    request: account_pb.ListRoleAssignmentsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.ListRoleAssignmentsResponse) => void
  ): grpcWeb.ClientReadableStream<account_pb.ListRoleAssignmentsResponse>;

  createRoleAssignment(
    request: account_pb.CreateRoleAssignmentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.RoleAssignment) => void
  ): grpcWeb.ClientReadableStream<account_pb.RoleAssignment>;

  deleteRoleAssignment(
    request: account_pb.DeleteRoleAssignmentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.DeleteRoleAssignmentResponse) => void
  ): grpcWeb.ClientReadableStream<account_pb.DeleteRoleAssignmentResponse>;

}

export class AccountInfoClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getPermission(
    request: account_pb.GetPermissionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.Permission) => void
  ): grpcWeb.ClientReadableStream<account_pb.Permission>;

  listPermissions(
    request: account_pb.ListPermissionsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.ListPermissionsResponse) => void
  ): grpcWeb.ClientReadableStream<account_pb.ListPermissionsResponse>;

  getAccountLimits(
    request: account_pb.GetAccountLimitsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: account_pb.AccountLimits) => void
  ): grpcWeb.ClientReadableStream<account_pb.AccountLimits>;

}

export class AccountApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getAccount(
    request: account_pb.GetAccountRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.Account>;

  listAccounts(
    request: account_pb.ListAccountsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.ListAccountsResponse>;

  createAccount(
    request: account_pb.CreateAccountRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.Account>;

  updateAccount(
    request: account_pb.UpdateAccountRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.Account>;

  updateAccountPassword(
    request: account_pb.UpdateAccountPasswordRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.UpdateAccountPasswordResponse>;

  rotateAccountClientSecret(
    request: account_pb.RotateAccountClientSecretRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.RotateAccountClientSecretResponse>;

  deleteAccount(
    request: account_pb.DeleteAccountRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.DeleteAccountResponse>;

  getRole(
    request: account_pb.GetRoleRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.Role>;

  listRoles(
    request: account_pb.ListRolesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.ListRolesResponse>;

  createRole(
    request: account_pb.CreateRoleRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.Role>;

  updateRole(
    request: account_pb.UpdateRoleRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.Role>;

  deleteRole(
    request: account_pb.DeleteRoleRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.DeleteRoleResponse>;

  getRoleAssignment(
    request: account_pb.GetRoleAssignmentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.RoleAssignment>;

  listRoleAssignments(
    request: account_pb.ListRoleAssignmentsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.ListRoleAssignmentsResponse>;

  createRoleAssignment(
    request: account_pb.CreateRoleAssignmentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.RoleAssignment>;

  deleteRoleAssignment(
    request: account_pb.DeleteRoleAssignmentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.DeleteRoleAssignmentResponse>;

}

export class AccountInfoPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getPermission(
    request: account_pb.GetPermissionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.Permission>;

  listPermissions(
    request: account_pb.ListPermissionsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.ListPermissionsResponse>;

  getAccountLimits(
    request: account_pb.GetAccountLimitsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<account_pb.AccountLimits>;

}

