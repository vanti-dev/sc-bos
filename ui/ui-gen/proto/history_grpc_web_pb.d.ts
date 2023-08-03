import * as grpcWeb from 'grpc-web';

import * as history_pb from './history_pb';


export class HistoryAdminApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  createHistoryRecord(
    request: history_pb.CreateHistoryRecordRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: history_pb.HistoryRecord) => void
  ): grpcWeb.ClientReadableStream<history_pb.HistoryRecord>;

  listHistoryRecords(
    request: history_pb.ListHistoryRecordsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: history_pb.ListHistoryRecordsResponse) => void
  ): grpcWeb.ClientReadableStream<history_pb.ListHistoryRecordsResponse>;

}

export class MeterHistoryClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listMeterReadingHistory(
    request: history_pb.ListMeterReadingHistoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: history_pb.ListMeterReadingHistoryResponse) => void
  ): grpcWeb.ClientReadableStream<history_pb.ListMeterReadingHistoryResponse>;

}

export class ElectricHistoryClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listElectricDemandHistory(
    request: history_pb.ListElectricDemandHistoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: history_pb.ListElectricDemandHistoryResponse) => void
  ): grpcWeb.ClientReadableStream<history_pb.ListElectricDemandHistoryResponse>;

}

export class OccupancySensorHistoryClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listOccupancyHistory(
    request: history_pb.ListOccupancyHistoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: history_pb.ListOccupancyHistoryResponse) => void
  ): grpcWeb.ClientReadableStream<history_pb.ListOccupancyHistoryResponse>;

}

export class HistoryAdminApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  createHistoryRecord(
    request: history_pb.CreateHistoryRecordRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<history_pb.HistoryRecord>;

  listHistoryRecords(
    request: history_pb.ListHistoryRecordsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<history_pb.ListHistoryRecordsResponse>;

}

export class MeterHistoryPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listMeterReadingHistory(
    request: history_pb.ListMeterReadingHistoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<history_pb.ListMeterReadingHistoryResponse>;

}

export class ElectricHistoryPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listElectricDemandHistory(
    request: history_pb.ListElectricDemandHistoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<history_pb.ListElectricDemandHistoryResponse>;

}

export class OccupancySensorHistoryPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listOccupancyHistory(
    request: history_pb.ListOccupancyHistoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<history_pb.ListOccupancyHistoryResponse>;

}

