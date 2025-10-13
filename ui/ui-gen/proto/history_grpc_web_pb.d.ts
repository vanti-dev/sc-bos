import * as grpcWeb from 'grpc-web';

import * as history_pb from './history_pb'; // proto import: "history.proto"


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

export class AirTemperatureHistoryClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listAirTemperatureHistory(
    request: history_pb.ListAirTemperatureHistoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: history_pb.ListAirTemperatureHistoryResponse) => void
  ): grpcWeb.ClientReadableStream<history_pb.ListAirTemperatureHistoryResponse>;

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

export class AirQualitySensorHistoryClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listAirQualityHistory(
    request: history_pb.ListAirQualityHistoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: history_pb.ListAirQualityHistoryResponse) => void
  ): grpcWeb.ClientReadableStream<history_pb.ListAirQualityHistoryResponse>;

}

export class SoundSensorHistoryClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listSoundLevelHistory(
    request: history_pb.ListSoundLevelHistoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: history_pb.ListSoundLevelHistoryResponse) => void
  ): grpcWeb.ClientReadableStream<history_pb.ListSoundLevelHistoryResponse>;

}

export class EnterLeaveHistoryClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listEnterLeaveSensorHistory(
    request: history_pb.ListEnterLeaveHistoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: history_pb.ListEnterLeaveHistoryResponse) => void
  ): grpcWeb.ClientReadableStream<history_pb.ListEnterLeaveHistoryResponse>;

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

export class AirTemperatureHistoryPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listAirTemperatureHistory(
    request: history_pb.ListAirTemperatureHistoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<history_pb.ListAirTemperatureHistoryResponse>;

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

export class AirQualitySensorHistoryPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listAirQualityHistory(
    request: history_pb.ListAirQualityHistoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<history_pb.ListAirQualityHistoryResponse>;

}

export class SoundSensorHistoryPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listSoundLevelHistory(
    request: history_pb.ListSoundLevelHistoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<history_pb.ListSoundLevelHistoryResponse>;

}

export class EnterLeaveHistoryPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listEnterLeaveSensorHistory(
    request: history_pb.ListEnterLeaveHistoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<history_pb.ListEnterLeaveHistoryResponse>;

}

