import * as grpcWeb from 'grpc-web';

import * as services_pb from './services_pb';


export class ServicesApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getService(
    request: services_pb.GetServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: services_pb.Service) => void
  ): grpcWeb.ClientReadableStream<services_pb.Service>;

  pullService(
    request: services_pb.PullServiceRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<services_pb.PullServiceResponse>;

  createService(
    request: services_pb.CreateServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: services_pb.Service) => void
  ): grpcWeb.ClientReadableStream<services_pb.Service>;

  deleteService(
    request: services_pb.DeleteServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: services_pb.Service) => void
  ): grpcWeb.ClientReadableStream<services_pb.Service>;

  listServices(
    request: services_pb.ListServicesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: services_pb.ListServicesResponse) => void
  ): grpcWeb.ClientReadableStream<services_pb.ListServicesResponse>;

  pullServices(
    request: services_pb.PullServicesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<services_pb.PullServicesResponse>;

  startService(
    request: services_pb.StartServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: services_pb.Service) => void
  ): grpcWeb.ClientReadableStream<services_pb.Service>;

  configureService(
    request: services_pb.ConfigureServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: services_pb.Service) => void
  ): grpcWeb.ClientReadableStream<services_pb.Service>;

  stopService(
    request: services_pb.StopServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: services_pb.Service) => void
  ): grpcWeb.ClientReadableStream<services_pb.Service>;

  getServiceMetadata(
    request: services_pb.GetServiceMetadataRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: services_pb.ServiceMetadata) => void
  ): grpcWeb.ClientReadableStream<services_pb.ServiceMetadata>;

  pullServiceMetadata(
    request: services_pb.PullServiceMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<services_pb.PullServiceMetadataResponse>;

}

export class ServicesApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getService(
    request: services_pb.GetServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<services_pb.Service>;

  pullService(
    request: services_pb.PullServiceRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<services_pb.PullServiceResponse>;

  createService(
    request: services_pb.CreateServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<services_pb.Service>;

  deleteService(
    request: services_pb.DeleteServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<services_pb.Service>;

  listServices(
    request: services_pb.ListServicesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<services_pb.ListServicesResponse>;

  pullServices(
    request: services_pb.PullServicesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<services_pb.PullServicesResponse>;

  startService(
    request: services_pb.StartServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<services_pb.Service>;

  configureService(
    request: services_pb.ConfigureServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<services_pb.Service>;

  stopService(
    request: services_pb.StopServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<services_pb.Service>;

  getServiceMetadata(
    request: services_pb.GetServiceMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<services_pb.ServiceMetadata>;

  pullServiceMetadata(
    request: services_pb.PullServiceMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<services_pb.PullServiceMetadataResponse>;

}

