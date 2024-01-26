import * as grpcWeb from 'grpc-web';

import * as meter_pb from './meter_pb'; // proto import: "meter.proto"


export class MeterApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getMeterReading(
    request: meter_pb.GetMeterReadingRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: meter_pb.MeterReading) => void
  ): grpcWeb.ClientReadableStream<meter_pb.MeterReading>;

  pullMeterReadings(
    request: meter_pb.PullMeterReadingsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<meter_pb.PullMeterReadingsResponse>;

}

export class MeterInfoClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describeMeterReading(
    request: meter_pb.DescribeMeterReadingRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: meter_pb.MeterReadingSupport) => void
  ): grpcWeb.ClientReadableStream<meter_pb.MeterReadingSupport>;

}

export class MeterApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getMeterReading(
    request: meter_pb.GetMeterReadingRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<meter_pb.MeterReading>;

  pullMeterReadings(
    request: meter_pb.PullMeterReadingsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<meter_pb.PullMeterReadingsResponse>;

}

export class MeterInfoPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describeMeterReading(
    request: meter_pb.DescribeMeterReadingRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<meter_pb.MeterReadingSupport>;

}

