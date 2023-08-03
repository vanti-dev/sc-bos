import * as grpcWeb from 'grpc-web';

import * as access_pb from './access_pb';


export class AccessApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getLastAccessAttempt(
    request: access_pb.GetLastAccessAttemptRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: access_pb.AccessAttempt) => void
  ): grpcWeb.ClientReadableStream<access_pb.AccessAttempt>;

  pullAccessAttempts(
    request: access_pb.PullAccessAttemptsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<access_pb.PullAccessAttemptsResponse>;

}

export class AccessApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getLastAccessAttempt(
    request: access_pb.GetLastAccessAttemptRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<access_pb.AccessAttempt>;

  pullAccessAttempts(
    request: access_pb.PullAccessAttemptsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<access_pb.PullAccessAttemptsResponse>;

}

