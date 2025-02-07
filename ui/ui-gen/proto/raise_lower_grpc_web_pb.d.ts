import * as grpcWeb from 'grpc-web';

import * as raise_lower_pb from './raise_lower_pb'; // proto import: "raise_lower.proto"


export class RaiseLowerApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getBearerState(
    request: raise_lower_pb.GetBearerStateRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: raise_lower_pb.BearerState) => void
  ): grpcWeb.ClientReadableStream<raise_lower_pb.BearerState>;

  pullBearerState(
    request: raise_lower_pb.PullBearerStateRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<raise_lower_pb.PullBearerStateResponse>;

}

export class RaiseLowerInfoClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describeBearerState(
    request: raise_lower_pb.DescribeBearerRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: raise_lower_pb.BearerStateSupport) => void
  ): grpcWeb.ClientReadableStream<raise_lower_pb.BearerStateSupport>;

}

export class RaiseLowerApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getBearerState(
    request: raise_lower_pb.GetBearerStateRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<raise_lower_pb.BearerState>;

  pullBearerState(
    request: raise_lower_pb.PullBearerStateRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<raise_lower_pb.PullBearerStateResponse>;

}

export class RaiseLowerInfoPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describeBearerState(
    request: raise_lower_pb.DescribeBearerRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<raise_lower_pb.BearerStateSupport>;

}

