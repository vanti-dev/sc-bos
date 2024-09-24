import * as grpcWeb from 'grpc-web';

import * as waste_pb from './waste_pb'; // proto import: "waste.proto"


export class WasteApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getWasteRecords(
    request: waste_pb.GetWasteRecordsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: waste_pb.WasteRecord) => void
  ): grpcWeb.ClientReadableStream<waste_pb.WasteRecord>;

  pullWasteRecords(
    request: waste_pb.PullWasteRecordsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<waste_pb.PullWasteRecordsResponse>;

}

export class WasteApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getWasteRecords(
    request: waste_pb.GetWasteRecordsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<waste_pb.WasteRecord>;

  pullWasteRecords(
    request: waste_pb.PullWasteRecordsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<waste_pb.PullWasteRecordsResponse>;

}

