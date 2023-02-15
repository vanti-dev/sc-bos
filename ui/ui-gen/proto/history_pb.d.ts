import * as jspb from 'google-protobuf'

import * as google_protobuf_field_mask_pb from 'google-protobuf/google/protobuf/field_mask_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as types_time_period_pb from '@smart-core-os/sc-api-grpc-web/types/time/period_pb';
import * as traits_electric_pb from '@smart-core-os/sc-api-grpc-web/traits/electric_pb';
import * as traits_occupancy_sensor_pb from '@smart-core-os/sc-api-grpc-web/traits/occupancy_sensor_pb';
import * as meter_pb from './meter_pb';


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

