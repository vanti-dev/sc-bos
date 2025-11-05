import * as grpcWeb from 'grpc-web';

import * as batch_data_pb from './batch_data_pb'; // proto import: "batch_data.proto"


export class BatchDataApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getBatchData(
    request: batch_data_pb.GetBatchDataRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: batch_data_pb.GetBatchDataResponse) => void
  ): grpcWeb.ClientReadableStream<batch_data_pb.GetBatchDataResponse>;

  pullBatchData(
    request: batch_data_pb.PullBatchDataRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<batch_data_pb.PullBatchDataResponse>;

}

export class BatchDataApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getBatchData(
    request: batch_data_pb.GetBatchDataRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<batch_data_pb.GetBatchDataResponse>;

  pullBatchData(
    request: batch_data_pb.PullBatchDataRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<batch_data_pb.PullBatchDataResponse>;

}

