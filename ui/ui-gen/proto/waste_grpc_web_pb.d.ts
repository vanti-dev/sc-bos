import * as grpcWeb from 'grpc-web';

import * as waste_pb from './waste_pb'; // proto import: "waste.proto"


export class WasteApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listWasteRecords(
    request: waste_pb.ListWasteRecordsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: waste_pb.ListWasteRecordsResponse) => void
  ): grpcWeb.ClientReadableStream<waste_pb.ListWasteRecordsResponse>;

  pullWasteRecords(
    request: waste_pb.PullWasteRecordsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<waste_pb.PullWasteRecordsResponse>;

}

export class WasteInfoClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describeWasteRecord(
    request: waste_pb.DescribeWasteRecordRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: waste_pb.WasteRecordSupport) => void
  ): grpcWeb.ClientReadableStream<waste_pb.WasteRecordSupport>;

}

export class WasteApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listWasteRecords(
    request: waste_pb.ListWasteRecordsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<waste_pb.ListWasteRecordsResponse>;

  pullWasteRecords(
    request: waste_pb.PullWasteRecordsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<waste_pb.PullWasteRecordsResponse>;

}

export class WasteInfoPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describeWasteRecord(
    request: waste_pb.DescribeWasteRecordRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<waste_pb.WasteRecordSupport>;

}

