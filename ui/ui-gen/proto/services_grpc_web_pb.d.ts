import * as grpcWeb from 'grpc-web';

import * as proto_services_pb from '../proto/services_pb';


export class ServicesApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getService(
    request: proto_services_pb.GetServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_services_pb.Service) => void
  ): grpcWeb.ClientReadableStream<proto_services_pb.Service>;

  pullService(
    request: proto_services_pb.PullServiceRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_services_pb.PullServiceResponse>;

  createService(
    request: proto_services_pb.CreateServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_services_pb.Service) => void
  ): grpcWeb.ClientReadableStream<proto_services_pb.Service>;

  deleteService(
    request: proto_services_pb.DeleteServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_services_pb.Service) => void
  ): grpcWeb.ClientReadableStream<proto_services_pb.Service>;

  listServices(
    request: proto_services_pb.ListServicesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_services_pb.ListServicesResponse) => void
  ): grpcWeb.ClientReadableStream<proto_services_pb.ListServicesResponse>;

  pullServices(
    request: proto_services_pb.PullServicesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_services_pb.PullServicesResponse>;

  startService(
    request: proto_services_pb.StartServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_services_pb.Service) => void
  ): grpcWeb.ClientReadableStream<proto_services_pb.Service>;

  configureService(
    request: proto_services_pb.ConfigureServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_services_pb.Service) => void
  ): grpcWeb.ClientReadableStream<proto_services_pb.Service>;

  stopService(
    request: proto_services_pb.StopServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_services_pb.Service) => void
  ): grpcWeb.ClientReadableStream<proto_services_pb.Service>;

  getServiceMetadata(
    request: proto_services_pb.GetServiceMetadataRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_services_pb.ServiceMetadata) => void
  ): grpcWeb.ClientReadableStream<proto_services_pb.ServiceMetadata>;

  pullServiceMetadata(
    request: proto_services_pb.PullServiceMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_services_pb.PullServiceMetadataResponse>;

}

export class ServicesApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getService(
    request: proto_services_pb.GetServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_services_pb.Service>;

  pullService(
    request: proto_services_pb.PullServiceRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_services_pb.PullServiceResponse>;

  createService(
    request: proto_services_pb.CreateServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_services_pb.Service>;

  deleteService(
    request: proto_services_pb.DeleteServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_services_pb.Service>;

  listServices(
    request: proto_services_pb.ListServicesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_services_pb.ListServicesResponse>;

  pullServices(
    request: proto_services_pb.PullServicesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_services_pb.PullServicesResponse>;

  startService(
    request: proto_services_pb.StartServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_services_pb.Service>;

  configureService(
    request: proto_services_pb.ConfigureServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_services_pb.Service>;

  stopService(
    request: proto_services_pb.StopServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_services_pb.Service>;

  getServiceMetadata(
    request: proto_services_pb.GetServiceMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_services_pb.ServiceMetadata>;

  pullServiceMetadata(
    request: proto_services_pb.PullServiceMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_services_pb.PullServiceMetadataResponse>;

}

