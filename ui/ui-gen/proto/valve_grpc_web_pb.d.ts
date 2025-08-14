import * as grpcWeb from 'grpc-web';

import * as valve_pb from './valve_pb'; // proto import: "valve.proto"


export class ValveApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getValve(
    request: valve_pb.GetValveRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: valve_pb.Valve) => void
  ): grpcWeb.ClientReadableStream<valve_pb.Valve>;

  pullValve(
    request: valve_pb.PullValveRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<valve_pb.PullValveResponse>;

  updateValve(
    request: valve_pb.UpdateValveRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: valve_pb.Valve) => void
  ): grpcWeb.ClientReadableStream<valve_pb.Valve>;

}

export class ValveApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getValve(
    request: valve_pb.GetValveRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<valve_pb.Valve>;

  pullValve(
    request: valve_pb.PullValveRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<valve_pb.PullValveResponse>;

  updateValve(
    request: valve_pb.UpdateValveRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<valve_pb.Valve>;

}

