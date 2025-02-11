import * as grpcWeb from 'grpc-web';

import * as transport_pb from './transport_pb'; // proto import: "transport.proto"


export class TransportApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getTransportState(
    request: transport_pb.GetTransportStateRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: transport_pb.TransportState) => void
  ): grpcWeb.ClientReadableStream<transport_pb.TransportState>;

  pullTransportState(
    request: transport_pb.PullTransportStateRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<transport_pb.PullTransportStateResponse>;

}

export class TransportInfoClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describeTransportState(
    request: transport_pb.DescribeTransportRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: transport_pb.TransportSupport) => void
  ): grpcWeb.ClientReadableStream<transport_pb.TransportSupport>;

}

export class TransportApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getTransportState(
    request: transport_pb.GetTransportStateRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<transport_pb.TransportState>;

  pullTransportState(
    request: transport_pb.PullTransportStateRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<transport_pb.PullTransportStateResponse>;

}

export class TransportInfoPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describeTransportState(
    request: transport_pb.DescribeTransportRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<transport_pb.TransportSupport>;

}

