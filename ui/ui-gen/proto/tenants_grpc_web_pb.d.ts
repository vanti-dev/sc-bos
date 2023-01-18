import * as grpcWeb from 'grpc-web';

import * as proto_tenants_pb from '../proto/tenants_pb';


export class TenantApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listTenants(
    request: proto_tenants_pb.ListTenantsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: proto_tenants_pb.ListTenantsResponse) => void
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.ListTenantsResponse>;

  pullTenants(
    request: proto_tenants_pb.PullTenantsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.PullTenantsResponse>;

  createTenant(
    request: proto_tenants_pb.CreateTenantRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: proto_tenants_pb.Tenant) => void
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.Tenant>;

  getTenant(
    request: proto_tenants_pb.GetTenantRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: proto_tenants_pb.Tenant) => void
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.Tenant>;

  updateTenant(
    request: proto_tenants_pb.UpdateTenantRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: proto_tenants_pb.Tenant) => void
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.Tenant>;

  deleteTenant(
    request: proto_tenants_pb.DeleteTenantRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: proto_tenants_pb.DeleteTenantResponse) => void
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.DeleteTenantResponse>;

  pullTenant(
    request: proto_tenants_pb.PullTenantRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.PullTenantResponse>;

  addTenantZones(
    request: proto_tenants_pb.AddTenantZonesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: proto_tenants_pb.Tenant) => void
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.Tenant>;

  removeTenantZones(
    request: proto_tenants_pb.RemoveTenantZonesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: proto_tenants_pb.Tenant) => void
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.Tenant>;

  listSecrets(
    request: proto_tenants_pb.ListSecretsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: proto_tenants_pb.ListSecretsResponse) => void
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.ListSecretsResponse>;

  pullSecrets(
    request: proto_tenants_pb.PullSecretsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.PullSecretsResponse>;

  createSecret(
    request: proto_tenants_pb.CreateSecretRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: proto_tenants_pb.Secret) => void
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.Secret>;

  verifySecret(
    request: proto_tenants_pb.VerifySecretRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: proto_tenants_pb.Secret) => void
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.Secret>;

  getSecret(
    request: proto_tenants_pb.GetSecretRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: proto_tenants_pb.Secret) => void
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.Secret>;

  updateSecret(
    request: proto_tenants_pb.UpdateSecretRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: proto_tenants_pb.Secret) => void
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.Secret>;

  deleteSecret(
    request: proto_tenants_pb.DeleteSecretRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: proto_tenants_pb.DeleteSecretResponse) => void
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.DeleteSecretResponse>;

  pullSecret(
    request: proto_tenants_pb.PullSecretRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.PullSecretResponse>;

  regenerateSecret(
    request: proto_tenants_pb.RegenerateSecretRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: proto_tenants_pb.Secret) => void
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.Secret>;

}

export class TenantApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listTenants(
    request: proto_tenants_pb.ListTenantsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_tenants_pb.ListTenantsResponse>;

  pullTenants(
    request: proto_tenants_pb.PullTenantsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.PullTenantsResponse>;

  createTenant(
    request: proto_tenants_pb.CreateTenantRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_tenants_pb.Tenant>;

  getTenant(
    request: proto_tenants_pb.GetTenantRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_tenants_pb.Tenant>;

  updateTenant(
    request: proto_tenants_pb.UpdateTenantRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_tenants_pb.Tenant>;

  deleteTenant(
    request: proto_tenants_pb.DeleteTenantRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_tenants_pb.DeleteTenantResponse>;

  pullTenant(
    request: proto_tenants_pb.PullTenantRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.PullTenantResponse>;

  addTenantZones(
    request: proto_tenants_pb.AddTenantZonesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_tenants_pb.Tenant>;

  removeTenantZones(
    request: proto_tenants_pb.RemoveTenantZonesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_tenants_pb.Tenant>;

  listSecrets(
    request: proto_tenants_pb.ListSecretsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_tenants_pb.ListSecretsResponse>;

  pullSecrets(
    request: proto_tenants_pb.PullSecretsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.PullSecretsResponse>;

  createSecret(
    request: proto_tenants_pb.CreateSecretRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_tenants_pb.Secret>;

  verifySecret(
    request: proto_tenants_pb.VerifySecretRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_tenants_pb.Secret>;

  getSecret(
    request: proto_tenants_pb.GetSecretRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_tenants_pb.Secret>;

  updateSecret(
    request: proto_tenants_pb.UpdateSecretRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_tenants_pb.Secret>;

  deleteSecret(
    request: proto_tenants_pb.DeleteSecretRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_tenants_pb.DeleteSecretResponse>;

  pullSecret(
    request: proto_tenants_pb.PullSecretRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_tenants_pb.PullSecretResponse>;

  regenerateSecret(
    request: proto_tenants_pb.RegenerateSecretRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_tenants_pb.Secret>;

}

