import * as jspb from 'google-protobuf'

import * as google_protobuf_field_mask_pb from 'google-protobuf/google/protobuf/field_mask_pb'; // proto import: "google/protobuf/field_mask.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as types_time_period_pb from '@smart-core-os/sc-api-grpc-web/types/time/period_pb'; // proto import: "types/time/period.proto"
import * as traits_air_temperature_pb from '@smart-core-os/sc-api-grpc-web/traits/air_temperature_pb'; // proto import: "traits/air_temperature.proto"
import * as traits_electric_pb from '@smart-core-os/sc-api-grpc-web/traits/electric_pb'; // proto import: "traits/electric.proto"
import * as traits_occupancy_sensor_pb from '@smart-core-os/sc-api-grpc-web/traits/occupancy_sensor_pb'; // proto import: "traits/occupancy_sensor.proto"
import * as traits_air_quality_sensor_pb from '@smart-core-os/sc-api-grpc-web/traits/air_quality_sensor_pb'; // proto import: "traits/air_quality_sensor.proto"
import * as meter_pb from './meter_pb'; // proto import: "meter.proto"


export class HistoryRecord extends jspb.Message {
  getId(): string;
  setId(value: string): HistoryRecord;

  getSource(): string;
  setSource(value: string): HistoryRecord;

  getCreateTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreateTime(value?: google_protobuf_timestamp_pb.Timestamp): HistoryRecord;
  hasCreateTime(): boolean;
  clearCreateTime(): HistoryRecord;

  getPayload(): Uint8Array | string;
  getPayload_asU8(): Uint8Array;
  getPayload_asB64(): string;
  setPayload(value: Uint8Array | string): HistoryRecord;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HistoryRecord.AsObject;
  static toObject(includeInstance: boolean, msg: HistoryRecord): HistoryRecord.AsObject;
  static serializeBinaryToWriter(message: HistoryRecord, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HistoryRecord;
  static deserializeBinaryFromReader(message: HistoryRecord, reader: jspb.BinaryReader): HistoryRecord;
}

export namespace HistoryRecord {
  export type AsObject = {
    id: string,
    source: string,
    createTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    payload: Uint8Array | string,
  }

  export class Query extends jspb.Message {
    getSourceEqual(): string;
    setSourceEqual(value: string): Query;

    getFromRecord(): HistoryRecord | undefined;
    setFromRecord(value?: HistoryRecord): Query;
    hasFromRecord(): boolean;
    clearFromRecord(): Query;

    getToRecord(): HistoryRecord | undefined;
    setToRecord(value?: HistoryRecord): Query;
    hasToRecord(): boolean;
    clearToRecord(): Query;

    getSourceCase(): Query.SourceCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Query.AsObject;
    static toObject(includeInstance: boolean, msg: Query): Query.AsObject;
    static serializeBinaryToWriter(message: Query, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Query;
    static deserializeBinaryFromReader(message: Query, reader: jspb.BinaryReader): Query;
  }

  export namespace Query {
    export type AsObject = {
      sourceEqual: string,
      fromRecord?: HistoryRecord.AsObject,
      toRecord?: HistoryRecord.AsObject,
    }

    export enum SourceCase { 
      SOURCE_NOT_SET = 0,
      SOURCE_EQUAL = 1,
    }
  }

}

export class CreateHistoryRecordRequest extends jspb.Message {
  getName(): string;
  setName(value: string): CreateHistoryRecordRequest;

  getRecord(): HistoryRecord | undefined;
  setRecord(value?: HistoryRecord): CreateHistoryRecordRequest;
  hasRecord(): boolean;
  clearRecord(): CreateHistoryRecordRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateHistoryRecordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateHistoryRecordRequest): CreateHistoryRecordRequest.AsObject;
  static serializeBinaryToWriter(message: CreateHistoryRecordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateHistoryRecordRequest;
  static deserializeBinaryFromReader(message: CreateHistoryRecordRequest, reader: jspb.BinaryReader): CreateHistoryRecordRequest;
}

export namespace CreateHistoryRecordRequest {
  export type AsObject = {
    name: string,
    record?: HistoryRecord.AsObject,
  }
}

export class ListHistoryRecordsRequest extends jspb.Message {
  getName(): string;
  setName(value: string): ListHistoryRecordsRequest;

  getPageSize(): number;
  setPageSize(value: number): ListHistoryRecordsRequest;

  getPageToken(): string;
  setPageToken(value: string): ListHistoryRecordsRequest;

  getOrderBy(): string;
  setOrderBy(value: string): ListHistoryRecordsRequest;

  getQuery(): HistoryRecord.Query | undefined;
  setQuery(value?: HistoryRecord.Query): ListHistoryRecordsRequest;
  hasQuery(): boolean;
  clearQuery(): ListHistoryRecordsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListHistoryRecordsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListHistoryRecordsRequest): ListHistoryRecordsRequest.AsObject;
  static serializeBinaryToWriter(message: ListHistoryRecordsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListHistoryRecordsRequest;
  static deserializeBinaryFromReader(message: ListHistoryRecordsRequest, reader: jspb.BinaryReader): ListHistoryRecordsRequest;
}

export namespace ListHistoryRecordsRequest {
  export type AsObject = {
    name: string,
    pageSize: number,
    pageToken: string,
    orderBy: string,
    query?: HistoryRecord.Query.AsObject,
  }
}

export class ListHistoryRecordsResponse extends jspb.Message {
  getRecordsList(): Array<HistoryRecord>;
  setRecordsList(value: Array<HistoryRecord>): ListHistoryRecordsResponse;
  clearRecordsList(): ListHistoryRecordsResponse;
  addRecords(value?: HistoryRecord, index?: number): HistoryRecord;

  getNextPageToken(): string;
  setNextPageToken(value: string): ListHistoryRecordsResponse;

  getTotalSize(): number;
  setTotalSize(value: number): ListHistoryRecordsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListHistoryRecordsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListHistoryRecordsResponse): ListHistoryRecordsResponse.AsObject;
  static serializeBinaryToWriter(message: ListHistoryRecordsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListHistoryRecordsResponse;
  static deserializeBinaryFromReader(message: ListHistoryRecordsResponse, reader: jspb.BinaryReader): ListHistoryRecordsResponse;
}

export namespace ListHistoryRecordsResponse {
  export type AsObject = {
    recordsList: Array<HistoryRecord.AsObject>,
    nextPageToken: string,
    totalSize: number,
  }
}

export class AirTemperatureRecord extends jspb.Message {
  getAirTemperature(): traits_air_temperature_pb.AirTemperature | undefined;
  setAirTemperature(value?: traits_air_temperature_pb.AirTemperature): AirTemperatureRecord;
  hasAirTemperature(): boolean;
  clearAirTemperature(): AirTemperatureRecord;

  getRecordTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setRecordTime(value?: google_protobuf_timestamp_pb.Timestamp): AirTemperatureRecord;
  hasRecordTime(): boolean;
  clearRecordTime(): AirTemperatureRecord;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AirTemperatureRecord.AsObject;
  static toObject(includeInstance: boolean, msg: AirTemperatureRecord): AirTemperatureRecord.AsObject;
  static serializeBinaryToWriter(message: AirTemperatureRecord, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AirTemperatureRecord;
  static deserializeBinaryFromReader(message: AirTemperatureRecord, reader: jspb.BinaryReader): AirTemperatureRecord;
}

export namespace AirTemperatureRecord {
  export type AsObject = {
    airTemperature?: traits_air_temperature_pb.AirTemperature.AsObject,
    recordTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class ListAirTemperatureHistoryRequest extends jspb.Message {
  getName(): string;
  setName(value: string): ListAirTemperatureHistoryRequest;

  getPeriod(): types_time_period_pb.Period | undefined;
  setPeriod(value?: types_time_period_pb.Period): ListAirTemperatureHistoryRequest;
  hasPeriod(): boolean;
  clearPeriod(): ListAirTemperatureHistoryRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): ListAirTemperatureHistoryRequest;
  hasReadMask(): boolean;
  clearReadMask(): ListAirTemperatureHistoryRequest;

  getPageSize(): number;
  setPageSize(value: number): ListAirTemperatureHistoryRequest;

  getPageToken(): string;
  setPageToken(value: string): ListAirTemperatureHistoryRequest;

  getOrderBy(): string;
  setOrderBy(value: string): ListAirTemperatureHistoryRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAirTemperatureHistoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAirTemperatureHistoryRequest): ListAirTemperatureHistoryRequest.AsObject;
  static serializeBinaryToWriter(message: ListAirTemperatureHistoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAirTemperatureHistoryRequest;
  static deserializeBinaryFromReader(message: ListAirTemperatureHistoryRequest, reader: jspb.BinaryReader): ListAirTemperatureHistoryRequest;
}

export namespace ListAirTemperatureHistoryRequest {
  export type AsObject = {
    name: string,
    period?: types_time_period_pb.Period.AsObject,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
    pageSize: number,
    pageToken: string,
    orderBy: string,
  }
}

export class ListAirTemperatureHistoryResponse extends jspb.Message {
  getAirTemperatureRecordsList(): Array<AirTemperatureRecord>;
  setAirTemperatureRecordsList(value: Array<AirTemperatureRecord>): ListAirTemperatureHistoryResponse;
  clearAirTemperatureRecordsList(): ListAirTemperatureHistoryResponse;
  addAirTemperatureRecords(value?: AirTemperatureRecord, index?: number): AirTemperatureRecord;

  getNextPageToken(): string;
  setNextPageToken(value: string): ListAirTemperatureHistoryResponse;

  getTotalSize(): number;
  setTotalSize(value: number): ListAirTemperatureHistoryResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAirTemperatureHistoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAirTemperatureHistoryResponse): ListAirTemperatureHistoryResponse.AsObject;
  static serializeBinaryToWriter(message: ListAirTemperatureHistoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAirTemperatureHistoryResponse;
  static deserializeBinaryFromReader(message: ListAirTemperatureHistoryResponse, reader: jspb.BinaryReader): ListAirTemperatureHistoryResponse;
}

export namespace ListAirTemperatureHistoryResponse {
  export type AsObject = {
    airTemperatureRecordsList: Array<AirTemperatureRecord.AsObject>,
    nextPageToken: string,
    totalSize: number,
  }
}

export class MeterReadingRecord extends jspb.Message {
  getMeterReading(): meter_pb.MeterReading | undefined;
  setMeterReading(value?: meter_pb.MeterReading): MeterReadingRecord;
  hasMeterReading(): boolean;
  clearMeterReading(): MeterReadingRecord;

  getRecordTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setRecordTime(value?: google_protobuf_timestamp_pb.Timestamp): MeterReadingRecord;
  hasRecordTime(): boolean;
  clearRecordTime(): MeterReadingRecord;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MeterReadingRecord.AsObject;
  static toObject(includeInstance: boolean, msg: MeterReadingRecord): MeterReadingRecord.AsObject;
  static serializeBinaryToWriter(message: MeterReadingRecord, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MeterReadingRecord;
  static deserializeBinaryFromReader(message: MeterReadingRecord, reader: jspb.BinaryReader): MeterReadingRecord;
}

export namespace MeterReadingRecord {
  export type AsObject = {
    meterReading?: meter_pb.MeterReading.AsObject,
    recordTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class ListMeterReadingHistoryRequest extends jspb.Message {
  getName(): string;
  setName(value: string): ListMeterReadingHistoryRequest;

  getPeriod(): types_time_period_pb.Period | undefined;
  setPeriod(value?: types_time_period_pb.Period): ListMeterReadingHistoryRequest;
  hasPeriod(): boolean;
  clearPeriod(): ListMeterReadingHistoryRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): ListMeterReadingHistoryRequest;
  hasReadMask(): boolean;
  clearReadMask(): ListMeterReadingHistoryRequest;

  getPageSize(): number;
  setPageSize(value: number): ListMeterReadingHistoryRequest;

  getPageToken(): string;
  setPageToken(value: string): ListMeterReadingHistoryRequest;

  getOrderBy(): string;
  setOrderBy(value: string): ListMeterReadingHistoryRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMeterReadingHistoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMeterReadingHistoryRequest): ListMeterReadingHistoryRequest.AsObject;
  static serializeBinaryToWriter(message: ListMeterReadingHistoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMeterReadingHistoryRequest;
  static deserializeBinaryFromReader(message: ListMeterReadingHistoryRequest, reader: jspb.BinaryReader): ListMeterReadingHistoryRequest;
}

export namespace ListMeterReadingHistoryRequest {
  export type AsObject = {
    name: string,
    period?: types_time_period_pb.Period.AsObject,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
    pageSize: number,
    pageToken: string,
    orderBy: string,
  }
}

export class ListMeterReadingHistoryResponse extends jspb.Message {
  getMeterReadingRecordsList(): Array<MeterReadingRecord>;
  setMeterReadingRecordsList(value: Array<MeterReadingRecord>): ListMeterReadingHistoryResponse;
  clearMeterReadingRecordsList(): ListMeterReadingHistoryResponse;
  addMeterReadingRecords(value?: MeterReadingRecord, index?: number): MeterReadingRecord;

  getNextPageToken(): string;
  setNextPageToken(value: string): ListMeterReadingHistoryResponse;

  getTotalSize(): number;
  setTotalSize(value: number): ListMeterReadingHistoryResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMeterReadingHistoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMeterReadingHistoryResponse): ListMeterReadingHistoryResponse.AsObject;
  static serializeBinaryToWriter(message: ListMeterReadingHistoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMeterReadingHistoryResponse;
  static deserializeBinaryFromReader(message: ListMeterReadingHistoryResponse, reader: jspb.BinaryReader): ListMeterReadingHistoryResponse;
}

export namespace ListMeterReadingHistoryResponse {
  export type AsObject = {
    meterReadingRecordsList: Array<MeterReadingRecord.AsObject>,
    nextPageToken: string,
    totalSize: number,
  }
}

export class ElectricDemandRecord extends jspb.Message {
  getElectricDemand(): traits_electric_pb.ElectricDemand | undefined;
  setElectricDemand(value?: traits_electric_pb.ElectricDemand): ElectricDemandRecord;
  hasElectricDemand(): boolean;
  clearElectricDemand(): ElectricDemandRecord;

  getRecordTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setRecordTime(value?: google_protobuf_timestamp_pb.Timestamp): ElectricDemandRecord;
  hasRecordTime(): boolean;
  clearRecordTime(): ElectricDemandRecord;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ElectricDemandRecord.AsObject;
  static toObject(includeInstance: boolean, msg: ElectricDemandRecord): ElectricDemandRecord.AsObject;
  static serializeBinaryToWriter(message: ElectricDemandRecord, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ElectricDemandRecord;
  static deserializeBinaryFromReader(message: ElectricDemandRecord, reader: jspb.BinaryReader): ElectricDemandRecord;
}

export namespace ElectricDemandRecord {
  export type AsObject = {
    electricDemand?: traits_electric_pb.ElectricDemand.AsObject,
    recordTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class ListElectricDemandHistoryRequest extends jspb.Message {
  getName(): string;
  setName(value: string): ListElectricDemandHistoryRequest;

  getPeriod(): types_time_period_pb.Period | undefined;
  setPeriod(value?: types_time_period_pb.Period): ListElectricDemandHistoryRequest;
  hasPeriod(): boolean;
  clearPeriod(): ListElectricDemandHistoryRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): ListElectricDemandHistoryRequest;
  hasReadMask(): boolean;
  clearReadMask(): ListElectricDemandHistoryRequest;

  getPageSize(): number;
  setPageSize(value: number): ListElectricDemandHistoryRequest;

  getPageToken(): string;
  setPageToken(value: string): ListElectricDemandHistoryRequest;

  getOrderBy(): string;
  setOrderBy(value: string): ListElectricDemandHistoryRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListElectricDemandHistoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListElectricDemandHistoryRequest): ListElectricDemandHistoryRequest.AsObject;
  static serializeBinaryToWriter(message: ListElectricDemandHistoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListElectricDemandHistoryRequest;
  static deserializeBinaryFromReader(message: ListElectricDemandHistoryRequest, reader: jspb.BinaryReader): ListElectricDemandHistoryRequest;
}

export namespace ListElectricDemandHistoryRequest {
  export type AsObject = {
    name: string,
    period?: types_time_period_pb.Period.AsObject,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
    pageSize: number,
    pageToken: string,
    orderBy: string,
  }
}

export class ListElectricDemandHistoryResponse extends jspb.Message {
  getElectricDemandRecordsList(): Array<ElectricDemandRecord>;
  setElectricDemandRecordsList(value: Array<ElectricDemandRecord>): ListElectricDemandHistoryResponse;
  clearElectricDemandRecordsList(): ListElectricDemandHistoryResponse;
  addElectricDemandRecords(value?: ElectricDemandRecord, index?: number): ElectricDemandRecord;

  getNextPageToken(): string;
  setNextPageToken(value: string): ListElectricDemandHistoryResponse;

  getTotalSize(): number;
  setTotalSize(value: number): ListElectricDemandHistoryResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListElectricDemandHistoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListElectricDemandHistoryResponse): ListElectricDemandHistoryResponse.AsObject;
  static serializeBinaryToWriter(message: ListElectricDemandHistoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListElectricDemandHistoryResponse;
  static deserializeBinaryFromReader(message: ListElectricDemandHistoryResponse, reader: jspb.BinaryReader): ListElectricDemandHistoryResponse;
}

export namespace ListElectricDemandHistoryResponse {
  export type AsObject = {
    electricDemandRecordsList: Array<ElectricDemandRecord.AsObject>,
    nextPageToken: string,
    totalSize: number,
  }
}

export class OccupancyRecord extends jspb.Message {
  getOccupancy(): traits_occupancy_sensor_pb.Occupancy | undefined;
  setOccupancy(value?: traits_occupancy_sensor_pb.Occupancy): OccupancyRecord;
  hasOccupancy(): boolean;
  clearOccupancy(): OccupancyRecord;

  getRecordTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setRecordTime(value?: google_protobuf_timestamp_pb.Timestamp): OccupancyRecord;
  hasRecordTime(): boolean;
  clearRecordTime(): OccupancyRecord;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OccupancyRecord.AsObject;
  static toObject(includeInstance: boolean, msg: OccupancyRecord): OccupancyRecord.AsObject;
  static serializeBinaryToWriter(message: OccupancyRecord, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OccupancyRecord;
  static deserializeBinaryFromReader(message: OccupancyRecord, reader: jspb.BinaryReader): OccupancyRecord;
}

export namespace OccupancyRecord {
  export type AsObject = {
    occupancy?: traits_occupancy_sensor_pb.Occupancy.AsObject,
    recordTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class ListOccupancyHistoryRequest extends jspb.Message {
  getName(): string;
  setName(value: string): ListOccupancyHistoryRequest;

  getPeriod(): types_time_period_pb.Period | undefined;
  setPeriod(value?: types_time_period_pb.Period): ListOccupancyHistoryRequest;
  hasPeriod(): boolean;
  clearPeriod(): ListOccupancyHistoryRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): ListOccupancyHistoryRequest;
  hasReadMask(): boolean;
  clearReadMask(): ListOccupancyHistoryRequest;

  getPageSize(): number;
  setPageSize(value: number): ListOccupancyHistoryRequest;

  getPageToken(): string;
  setPageToken(value: string): ListOccupancyHistoryRequest;

  getOrderBy(): string;
  setOrderBy(value: string): ListOccupancyHistoryRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOccupancyHistoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListOccupancyHistoryRequest): ListOccupancyHistoryRequest.AsObject;
  static serializeBinaryToWriter(message: ListOccupancyHistoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOccupancyHistoryRequest;
  static deserializeBinaryFromReader(message: ListOccupancyHistoryRequest, reader: jspb.BinaryReader): ListOccupancyHistoryRequest;
}

export namespace ListOccupancyHistoryRequest {
  export type AsObject = {
    name: string,
    period?: types_time_period_pb.Period.AsObject,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
    pageSize: number,
    pageToken: string,
    orderBy: string,
  }
}

export class ListOccupancyHistoryResponse extends jspb.Message {
  getOccupancyRecordsList(): Array<OccupancyRecord>;
  setOccupancyRecordsList(value: Array<OccupancyRecord>): ListOccupancyHistoryResponse;
  clearOccupancyRecordsList(): ListOccupancyHistoryResponse;
  addOccupancyRecords(value?: OccupancyRecord, index?: number): OccupancyRecord;

  getNextPageToken(): string;
  setNextPageToken(value: string): ListOccupancyHistoryResponse;

  getTotalSize(): number;
  setTotalSize(value: number): ListOccupancyHistoryResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOccupancyHistoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListOccupancyHistoryResponse): ListOccupancyHistoryResponse.AsObject;
  static serializeBinaryToWriter(message: ListOccupancyHistoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOccupancyHistoryResponse;
  static deserializeBinaryFromReader(message: ListOccupancyHistoryResponse, reader: jspb.BinaryReader): ListOccupancyHistoryResponse;
}

export namespace ListOccupancyHistoryResponse {
  export type AsObject = {
    occupancyRecordsList: Array<OccupancyRecord.AsObject>,
    nextPageToken: string,
    totalSize: number,
  }
}

export class AirQualityRecord extends jspb.Message {
  getAirQuality(): traits_air_quality_sensor_pb.AirQuality | undefined;
  setAirQuality(value?: traits_air_quality_sensor_pb.AirQuality): AirQualityRecord;
  hasAirQuality(): boolean;
  clearAirQuality(): AirQualityRecord;

  getRecordTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setRecordTime(value?: google_protobuf_timestamp_pb.Timestamp): AirQualityRecord;
  hasRecordTime(): boolean;
  clearRecordTime(): AirQualityRecord;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AirQualityRecord.AsObject;
  static toObject(includeInstance: boolean, msg: AirQualityRecord): AirQualityRecord.AsObject;
  static serializeBinaryToWriter(message: AirQualityRecord, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AirQualityRecord;
  static deserializeBinaryFromReader(message: AirQualityRecord, reader: jspb.BinaryReader): AirQualityRecord;
}

export namespace AirQualityRecord {
  export type AsObject = {
    airQuality?: traits_air_quality_sensor_pb.AirQuality.AsObject,
    recordTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class ListAirQualityHistoryRequest extends jspb.Message {
  getName(): string;
  setName(value: string): ListAirQualityHistoryRequest;

  getPeriod(): types_time_period_pb.Period | undefined;
  setPeriod(value?: types_time_period_pb.Period): ListAirQualityHistoryRequest;
  hasPeriod(): boolean;
  clearPeriod(): ListAirQualityHistoryRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): ListAirQualityHistoryRequest;
  hasReadMask(): boolean;
  clearReadMask(): ListAirQualityHistoryRequest;

  getPageSize(): number;
  setPageSize(value: number): ListAirQualityHistoryRequest;

  getPageToken(): string;
  setPageToken(value: string): ListAirQualityHistoryRequest;

  getOrderBy(): string;
  setOrderBy(value: string): ListAirQualityHistoryRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAirQualityHistoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAirQualityHistoryRequest): ListAirQualityHistoryRequest.AsObject;
  static serializeBinaryToWriter(message: ListAirQualityHistoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAirQualityHistoryRequest;
  static deserializeBinaryFromReader(message: ListAirQualityHistoryRequest, reader: jspb.BinaryReader): ListAirQualityHistoryRequest;
}

export namespace ListAirQualityHistoryRequest {
  export type AsObject = {
    name: string,
    period?: types_time_period_pb.Period.AsObject,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
    pageSize: number,
    pageToken: string,
    orderBy: string,
  }
}

export class ListAirQualityHistoryResponse extends jspb.Message {
  getAirQualityRecordsList(): Array<AirQualityRecord>;
  setAirQualityRecordsList(value: Array<AirQualityRecord>): ListAirQualityHistoryResponse;
  clearAirQualityRecordsList(): ListAirQualityHistoryResponse;
  addAirQualityRecords(value?: AirQualityRecord, index?: number): AirQualityRecord;

  getNextPageToken(): string;
  setNextPageToken(value: string): ListAirQualityHistoryResponse;

  getTotalSize(): number;
  setTotalSize(value: number): ListAirQualityHistoryResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAirQualityHistoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAirQualityHistoryResponse): ListAirQualityHistoryResponse.AsObject;
  static serializeBinaryToWriter(message: ListAirQualityHistoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAirQualityHistoryResponse;
  static deserializeBinaryFromReader(message: ListAirQualityHistoryResponse, reader: jspb.BinaryReader): ListAirQualityHistoryResponse;
}

export namespace ListAirQualityHistoryResponse {
  export type AsObject = {
    airQualityRecordsList: Array<AirQualityRecord.AsObject>,
    nextPageToken: string,
    totalSize: number,
  }
}

