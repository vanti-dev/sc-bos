import * as grpcWeb from 'grpc-web';

import * as temperature_pb from './temperature_pb'; // proto import: "temperature.proto"


export class TemperatureApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getTemperature(
    request: temperature_pb.GetTemperatureRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: temperature_pb.GetTemperatureResponse) => void
  ): grpcWeb.ClientReadableStream<temperature_pb.GetTemperatureResponse>;

  pullTemperature(
    request: temperature_pb.PullTemperatureRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<temperature_pb.PullTemperatureResponse>;

  updateTemperature(
    request: temperature_pb.UpdateTemperatureRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: temperature_pb.UpdateTemperatureResponse) => void
  ): grpcWeb.ClientReadableStream<temperature_pb.UpdateTemperatureResponse>;

}

export class TemperatureApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getTemperature(
    request: temperature_pb.GetTemperatureRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<temperature_pb.GetTemperatureResponse>;

  pullTemperature(
    request: temperature_pb.PullTemperatureRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<temperature_pb.PullTemperatureResponse>;

  updateTemperature(
    request: temperature_pb.UpdateTemperatureRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<temperature_pb.UpdateTemperatureResponse>;

}

