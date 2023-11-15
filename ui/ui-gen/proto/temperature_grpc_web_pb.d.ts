import * as grpcWeb from 'grpc-web';

import * as temperature_pb from './temperature_pb';


export class TemperatureApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getTemperature(
    request: temperature_pb.GetTemperatureRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: temperature_pb.Temperature) => void
  ): grpcWeb.ClientReadableStream<temperature_pb.Temperature>;

  pullTemperature(
    request: temperature_pb.PullTemperatureRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<temperature_pb.PullTemperatureResponse>;

  updateTemperature(
    request: temperature_pb.UpdateTemperatureRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: temperature_pb.Temperature) => void
  ): grpcWeb.ClientReadableStream<temperature_pb.Temperature>;

}

export class TemperatureApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getTemperature(
    request: temperature_pb.GetTemperatureRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<temperature_pb.Temperature>;

  pullTemperature(
    request: temperature_pb.PullTemperatureRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<temperature_pb.PullTemperatureResponse>;

  updateTemperature(
    request: temperature_pb.UpdateTemperatureRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<temperature_pb.Temperature>;

}

