import * as jspb from 'google-protobuf'

import * as google_protobuf_any_pb from 'google-protobuf/google/protobuf/any_pb'; // proto import: "google/protobuf/any.proto"
import * as google_protobuf_field_mask_pb from 'google-protobuf/google/protobuf/field_mask_pb'; // proto import: "google/protobuf/field_mask.proto"


export class BatchData extends jspb.Message {
  getName(): string;
  setName(value: string): BatchData;

  getTrait(): string;
  setTrait(value: string): BatchData;

  getData(): google_protobuf_any_pb.Any | undefined;
  setData(value?: google_protobuf_any_pb.Any): BatchData;
  hasData(): boolean;
  clearData(): BatchData;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BatchData.AsObject;
  static toObject(includeInstance: boolean, msg: BatchData): BatchData.AsObject;
  static serializeBinaryToWriter(message: BatchData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BatchData;
  static deserializeBinaryFromReader(message: BatchData, reader: jspb.BinaryReader): BatchData;
}

export namespace BatchData {
  export type AsObject = {
    name: string,
    trait: string,
    data?: google_protobuf_any_pb.Any.AsObject,
  }
}

export class GetBatchDataRequest extends jspb.Message {
  getSourcesList(): Array<BatchDataSource>;
  setSourcesList(value: Array<BatchDataSource>): GetBatchDataRequest;
  clearSourcesList(): GetBatchDataRequest;
  addSources(value?: BatchDataSource, index?: number): BatchDataSource;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetBatchDataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetBatchDataRequest): GetBatchDataRequest.AsObject;
  static serializeBinaryToWriter(message: GetBatchDataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetBatchDataRequest;
  static deserializeBinaryFromReader(message: GetBatchDataRequest, reader: jspb.BinaryReader): GetBatchDataRequest;
}

export namespace GetBatchDataRequest {
  export type AsObject = {
    sourcesList: Array<BatchDataSource.AsObject>,
  }
}

export class BatchDataSource extends jspb.Message {
  getName(): string;
  setName(value: string): BatchDataSource;

  getTrait(): string;
  setTrait(value: string): BatchDataSource;

  getFieldMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setFieldMask(value?: google_protobuf_field_mask_pb.FieldMask): BatchDataSource;
  hasFieldMask(): boolean;
  clearFieldMask(): BatchDataSource;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BatchDataSource.AsObject;
  static toObject(includeInstance: boolean, msg: BatchDataSource): BatchDataSource.AsObject;
  static serializeBinaryToWriter(message: BatchDataSource, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BatchDataSource;
  static deserializeBinaryFromReader(message: BatchDataSource, reader: jspb.BinaryReader): BatchDataSource;
}

export namespace BatchDataSource {
  export type AsObject = {
    name: string,
    trait: string,
    fieldMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
  }
}

export class GetBatchDataResponse extends jspb.Message {
  getDataList(): Array<BatchData>;
  setDataList(value: Array<BatchData>): GetBatchDataResponse;
  clearDataList(): GetBatchDataResponse;
  addData(value?: BatchData, index?: number): BatchData;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetBatchDataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetBatchDataResponse): GetBatchDataResponse.AsObject;
  static serializeBinaryToWriter(message: GetBatchDataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetBatchDataResponse;
  static deserializeBinaryFromReader(message: GetBatchDataResponse, reader: jspb.BinaryReader): GetBatchDataResponse;
}

export namespace GetBatchDataResponse {
  export type AsObject = {
    dataList: Array<BatchData.AsObject>,
  }
}

export class PullBatchDataRequest extends jspb.Message {
  getSourcesList(): Array<BatchDataSource>;
  setSourcesList(value: Array<BatchDataSource>): PullBatchDataRequest;
  clearSourcesList(): PullBatchDataRequest;
  addSources(value?: BatchDataSource, index?: number): BatchDataSource;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullBatchDataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PullBatchDataRequest): PullBatchDataRequest.AsObject;
  static serializeBinaryToWriter(message: PullBatchDataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullBatchDataRequest;
  static deserializeBinaryFromReader(message: PullBatchDataRequest, reader: jspb.BinaryReader): PullBatchDataRequest;
}

export namespace PullBatchDataRequest {
  export type AsObject = {
    sourcesList: Array<BatchDataSource.AsObject>,
  }
}

export class PullBatchDataResponse extends jspb.Message {
  getDataList(): Array<BatchData>;
  setDataList(value: Array<BatchData>): PullBatchDataResponse;
  clearDataList(): PullBatchDataResponse;
  addData(value?: BatchData, index?: number): BatchData;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullBatchDataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PullBatchDataResponse): PullBatchDataResponse.AsObject;
  static serializeBinaryToWriter(message: PullBatchDataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullBatchDataResponse;
  static deserializeBinaryFromReader(message: PullBatchDataResponse, reader: jspb.BinaryReader): PullBatchDataResponse;
}

export namespace PullBatchDataResponse {
  export type AsObject = {
    dataList: Array<BatchData.AsObject>,
  }
}

