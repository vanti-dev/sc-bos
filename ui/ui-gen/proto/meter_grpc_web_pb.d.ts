import * as grpcWeb from 'grpc-web';

import * as proto_meter_pb from '../proto/meter_pb';


export class MeterApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getMeterReading(
    request: proto_meter_pb.GetMeterReadingRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_meter_pb.MeterReading) => void
  ): grpcWeb.ClientReadableStream<proto_meter_pb.MeterReading>;

  pullMeterReadings(
    request: proto_meter_pb.PullMeterReadingsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_meter_pb.PullMeterReadingsResponse>;

}

export class MeterApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getMeterReading(
    request: proto_meter_pb.GetMeterReadingRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_meter_pb.MeterReading>;

  pullMeterReadings(
    request: proto_meter_pb.PullMeterReadingsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_meter_pb.PullMeterReadingsResponse>;

}

