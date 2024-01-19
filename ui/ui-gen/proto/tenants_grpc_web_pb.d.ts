import * as grpcWeb from 'grpc-web';

import * as tenants_pb from './tenants_pb'; // proto import: "tenants.proto"


export class TenantApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listTenants(
    request: tenants_pb.ListTenantsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: tenants_pb.ListTenantsResponse) => void
  ): grpcWeb.ClientReadableStream<tenants_pb.ListTenantsResponse>;

  pullTenants(
    request: tenants_pb.PullTenantsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<tenants_pb.PullTenantsResponse>;

  createTenant(
    request: tenants_pb.CreateTenantRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: tenants_pb.Tenant) => void
  ): grpcWeb.ClientReadableStream<tenants_pb.Tenant>;

  getTenant(
    request: tenants_pb.GetTenantRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: tenants_pb.Tenant) => void
  ): grpcWeb.ClientReadableStream<tenants_pb.Tenant>;

  updateTenant(
    request: tenants_pb.UpdateTenantRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: tenants_pb.Tenant) => void
  ): grpcWeb.ClientReadableStream<tenants_pb.Tenant>;

  deleteTenant(
    request: tenants_pb.DeleteTenantRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: tenants_pb.DeleteTenantResponse) => void
  ): grpcWeb.ClientReadableStream<tenants_pb.DeleteTenantResponse>;

  pullTenant(
    request: tenants_pb.PullTenantRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<tenants_pb.PullTenantResponse>;

  addTenantZones(
    request: tenants_pb.AddTenantZonesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: tenants_pb.Tenant) => void
  ): grpcWeb.ClientReadableStream<tenants_pb.Tenant>;

  removeTenantZones(
    request: tenants_pb.RemoveTenantZonesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: tenants_pb.Tenant) => void
  ): grpcWeb.ClientReadableStream<tenants_pb.Tenant>;

  listSecrets(
    request: tenants_pb.ListSecretsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: tenants_pb.ListSecretsResponse) => void
  ): grpcWeb.ClientReadableStream<tenants_pb.ListSecretsResponse>;

  pullSecrets(
    request: tenants_pb.PullSecretsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<tenants_pb.PullSecretsResponse>;

  createSecret(
    request: tenants_pb.CreateSecretRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: tenants_pb.Secret) => void
  ): grpcWeb.ClientReadableStream<tenants_pb.Secret>;

  verifySecret(
    request: tenants_pb.VerifySecretRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: tenants_pb.Secret) => void
  ): grpcWeb.ClientReadableStream<tenants_pb.Secret>;

  getSecret(
    request: tenants_pb.GetSecretRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: tenants_pb.Secret) => void
  ): grpcWeb.ClientReadableStream<tenants_pb.Secret>;

  updateSecret(
    request: tenants_pb.UpdateSecretRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: tenants_pb.Secret) => void
  ): grpcWeb.ClientReadableStream<tenants_pb.Secret>;

  deleteSecret(
    request: tenants_pb.DeleteSecretRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: tenants_pb.DeleteSecretResponse) => void
  ): grpcWeb.ClientReadableStream<tenants_pb.DeleteSecretResponse>;

  pullSecret(
    request: tenants_pb.PullSecretRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<tenants_pb.PullSecretResponse>;

  regenerateSecret(
    request: tenants_pb.RegenerateSecretRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: tenants_pb.Secret) => void
  ): grpcWeb.ClientReadableStream<tenants_pb.Secret>;

}

export class TenantApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listTenants(
    request: tenants_pb.ListTenantsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<tenants_pb.ListTenantsResponse>;

  pullTenants(
    request: tenants_pb.PullTenantsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<tenants_pb.PullTenantsResponse>;

  createTenant(
    request: tenants_pb.CreateTenantRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<tenants_pb.Tenant>;

  getTenant(
    request: tenants_pb.GetTenantRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<tenants_pb.Tenant>;

  updateTenant(
    request: tenants_pb.UpdateTenantRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<tenants_pb.Tenant>;

  deleteTenant(
    request: tenants_pb.DeleteTenantRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<tenants_pb.DeleteTenantResponse>;

  pullTenant(
    request: tenants_pb.PullTenantRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<tenants_pb.PullTenantResponse>;

  addTenantZones(
    request: tenants_pb.AddTenantZonesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<tenants_pb.Tenant>;

  removeTenantZones(
    request: tenants_pb.RemoveTenantZonesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<tenants_pb.Tenant>;

  listSecrets(
    request: tenants_pb.ListSecretsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<tenants_pb.ListSecretsResponse>;

  pullSecrets(
    request: tenants_pb.PullSecretsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<tenants_pb.PullSecretsResponse>;

  createSecret(
    request: tenants_pb.CreateSecretRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<tenants_pb.Secret>;

  verifySecret(
    request: tenants_pb.VerifySecretRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<tenants_pb.Secret>;

  getSecret(
    request: tenants_pb.GetSecretRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<tenants_pb.Secret>;

  updateSecret(
    request: tenants_pb.UpdateSecretRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<tenants_pb.Secret>;

  deleteSecret(
    request: tenants_pb.DeleteSecretRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<tenants_pb.DeleteSecretResponse>;

  pullSecret(
    request: tenants_pb.PullSecretRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<tenants_pb.PullSecretResponse>;

  regenerateSecret(
    request: tenants_pb.RegenerateSecretRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<tenants_pb.Secret>;

}

