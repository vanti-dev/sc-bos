import * as grpcWeb from 'grpc-web';

import * as proto_test_pb from '../proto/test_pb';


export class TestApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getTest(
    request: proto_test_pb.GetTestRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_test_pb.Test) => void
  ): grpcWeb.ClientReadableStream<proto_test_pb.Test>;

  updateTest(
    request: proto_test_pb.UpdateTestRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_test_pb.Test) => void
  ): grpcWeb.ClientReadableStream<proto_test_pb.Test>;

}

export class TestApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getTest(
    request: proto_test_pb.GetTestRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_test_pb.Test>;

  updateTest(
    request: proto_test_pb.UpdateTestRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_test_pb.Test>;

}

