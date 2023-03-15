import * as grpcWeb from 'grpc-web';

import * as color_pb from './color_pb';


export class ColorApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getColor(
    request: color_pb.GetColorRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: color_pb.Color) => void
  ): grpcWeb.ClientReadableStream<color_pb.Color>;

  updateColor(
    request: color_pb.UpdateColorRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: color_pb.Color) => void
  ): grpcWeb.ClientReadableStream<color_pb.Color>;

  pullColor(
    request: color_pb.PullColorRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<color_pb.PullColorResponse>;

}

export class ColorInfoClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describeColor(
    request: color_pb.DescribeColorRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: color_pb.ColorSupport) => void
  ): grpcWeb.ClientReadableStream<color_pb.ColorSupport>;

}

export class ColorApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getColor(
    request: color_pb.GetColorRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<color_pb.Color>;

  updateColor(
    request: color_pb.UpdateColorRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<color_pb.Color>;

  pullColor(
    request: color_pb.PullColorRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<color_pb.PullColorResponse>;

}

export class ColorInfoPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describeColor(
    request: color_pb.DescribeColorRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<color_pb.ColorSupport>;

}

