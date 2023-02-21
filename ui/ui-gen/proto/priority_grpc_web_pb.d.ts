import * as grpcWeb from 'grpc-web';

import * as priority_pb from './priority_pb';


export class PriorityApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  clearPriorityEntry(
    request: priority_pb.ClearPriorityValueRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: priority_pb.ClearPriorityValueResponse) => void
  ): grpcWeb.ClientReadableStream<priority_pb.ClearPriorityValueResponse>;

}

export class PriorityApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  clearPriorityEntry(
    request: priority_pb.ClearPriorityValueRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<priority_pb.ClearPriorityValueResponse>;

}

