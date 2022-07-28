import * as grpcWeb from 'grpc-web';

import * as test_pb from './test_pb';


export class TestApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getTest(
    request: test_pb.GetTestRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: test_pb.Test) => void
  ): grpcWeb.ClientReadableStream<test_pb.Test>;

  updateTest(
    request: test_pb.UpdateTestRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: test_pb.Test) => void
  ): grpcWeb.ClientReadableStream<test_pb.Test>;

}

export class TestApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getTest(
    request: test_pb.GetTestRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<test_pb.Test>;

  updateTest(
    request: test_pb.UpdateTestRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<test_pb.Test>;

}

