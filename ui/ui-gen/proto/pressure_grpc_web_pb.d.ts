import * as grpcWeb from 'grpc-web';

import * as pressure_pb from './pressure_pb'; // proto import: "pressure.proto"


export class PressureApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getPressure(
    request: pressure_pb.GetPressureRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pressure_pb.Pressure) => void
  ): grpcWeb.ClientReadableStream<pressure_pb.Pressure>;

  pullPressure(
    request: pressure_pb.PullPressureRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<pressure_pb.PullPressureResponse>;

  updatePressure(
    request: pressure_pb.UpdatePressureRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pressure_pb.Pressure) => void
  ): grpcWeb.ClientReadableStream<pressure_pb.Pressure>;

}

export class PressureInfoClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describePressure(
    request: pressure_pb.DescribePressureRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pressure_pb.PressureSupport) => void
  ): grpcWeb.ClientReadableStream<pressure_pb.PressureSupport>;

}

export class PressureApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getPressure(
    request: pressure_pb.GetPressureRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pressure_pb.Pressure>;

  pullPressure(
    request: pressure_pb.PullPressureRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<pressure_pb.PullPressureResponse>;

  updatePressure(
    request: pressure_pb.UpdatePressureRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pressure_pb.Pressure>;

}

export class PressureInfoPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describePressure(
    request: pressure_pb.DescribePressureRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pressure_pb.PressureSupport>;

}

