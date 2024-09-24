import * as jspb from 'google-protobuf'

import * as google_protobuf_field_mask_pb from 'google-protobuf/google/protobuf/field_mask_pb'; // proto import: "google/protobuf/field_mask.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as types_info_pb from '@smart-core-os/sc-api-grpc-web/types/info_pb'; // proto import: "types/info.proto"


export class WasteRecord extends jspb.Message {
  getId(): string;
  setId(value: string): WasteRecord;

  getRecordcreatedtime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setRecordcreatedtime(value?: google_protobuf_timestamp_pb.Timestamp): WasteRecord;
  hasRecordcreatedtime(): boolean;
  clearRecordcreatedtime(): WasteRecord;

  getWeight(): number;
  setWeight(value: number): WasteRecord;

  getSystem(): string;
  setSystem(value: string): WasteRecord;

  getRecycled(): boolean;
  setRecycled(value: boolean): WasteRecord;

  getArea(): string;
  setArea(value: string): WasteRecord;

  getWastecreateddate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setWastecreateddate(value?: google_protobuf_timestamp_pb.Timestamp): WasteRecord;
  hasWastecreateddate(): boolean;
  clearWastecreateddate(): WasteRecord;

  getStream(): string;
  setStream(value: string): WasteRecord;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WasteRecord.AsObject;
  static toObject(includeInstance: boolean, msg: WasteRecord): WasteRecord.AsObject;
  static serializeBinaryToWriter(message: WasteRecord, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WasteRecord;
  static deserializeBinaryFromReader(message: WasteRecord, reader: jspb.BinaryReader): WasteRecord;
}

export namespace WasteRecord {
  export type AsObject = {
    id: string,
    recordcreatedtime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    weight: number,
    system: string,
    recycled: boolean,
    area: string,
    wastecreateddate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    stream: string,
  }
}

export class ListWasteRecordsResponse extends jspb.Message {
  getWasterecordsList(): Array<WasteRecord>;
  setWasterecordsList(value: Array<WasteRecord>): ListWasteRecordsResponse;
  clearWasterecordsList(): ListWasteRecordsResponse;
  addWasterecords(value?: WasteRecord, index?: number): WasteRecord;

  getNextPageToken(): string;
  setNextPageToken(value: string): ListWasteRecordsResponse;

  getTotalSize(): number;
  setTotalSize(value: number): ListWasteRecordsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListWasteRecordsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListWasteRecordsResponse): ListWasteRecordsResponse.AsObject;
  static serializeBinaryToWriter(message: ListWasteRecordsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListWasteRecordsResponse;
  static deserializeBinaryFromReader(message: ListWasteRecordsResponse, reader: jspb.BinaryReader): ListWasteRecordsResponse;
}

export namespace ListWasteRecordsResponse {
  export type AsObject = {
    wasterecordsList: Array<WasteRecord.AsObject>,
    nextPageToken: string,
    totalSize: number,
  }
}

export class ListWasteRecordsRequest extends jspb.Message {
  getName(): string;
  setName(value: string): ListWasteRecordsRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): ListWasteRecordsRequest;
  hasReadMask(): boolean;
  clearReadMask(): ListWasteRecordsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListWasteRecordsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListWasteRecordsRequest): ListWasteRecordsRequest.AsObject;
  static serializeBinaryToWriter(message: ListWasteRecordsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListWasteRecordsRequest;
  static deserializeBinaryFromReader(message: ListWasteRecordsRequest, reader: jspb.BinaryReader): ListWasteRecordsRequest;
}

export namespace ListWasteRecordsRequest {
  export type AsObject = {
    name: string,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
  }
}

export class PullWasteRecordsRequest extends jspb.Message {
  getName(): string;
  setName(value: string): PullWasteRecordsRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): PullWasteRecordsRequest;
  hasReadMask(): boolean;
  clearReadMask(): PullWasteRecordsRequest;

  getUpdatesOnly(): boolean;
  setUpdatesOnly(value: boolean): PullWasteRecordsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullWasteRecordsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PullWasteRecordsRequest): PullWasteRecordsRequest.AsObject;
  static serializeBinaryToWriter(message: PullWasteRecordsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullWasteRecordsRequest;
  static deserializeBinaryFromReader(message: PullWasteRecordsRequest, reader: jspb.BinaryReader): PullWasteRecordsRequest;
}

export namespace PullWasteRecordsRequest {
  export type AsObject = {
    name: string,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
    updatesOnly: boolean,
  }
}

export class PullWasteRecordsResponse extends jspb.Message {
  getChangesList(): Array<PullWasteRecordsResponse.Change>;
  setChangesList(value: Array<PullWasteRecordsResponse.Change>): PullWasteRecordsResponse;
  clearChangesList(): PullWasteRecordsResponse;
  addChanges(value?: PullWasteRecordsResponse.Change, index?: number): PullWasteRecordsResponse.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullWasteRecordsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PullWasteRecordsResponse): PullWasteRecordsResponse.AsObject;
  static serializeBinaryToWriter(message: PullWasteRecordsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullWasteRecordsResponse;
  static deserializeBinaryFromReader(message: PullWasteRecordsResponse, reader: jspb.BinaryReader): PullWasteRecordsResponse;
}

export namespace PullWasteRecordsResponse {
  export type AsObject = {
    changesList: Array<PullWasteRecordsResponse.Change.AsObject>,
  }

  export class Change extends jspb.Message {
    getName(): string;
    setName(value: string): Change;

    getChangeTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setChangeTime(value?: google_protobuf_timestamp_pb.Timestamp): Change;
    hasChangeTime(): boolean;
    clearChangeTime(): Change;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Change.AsObject;
    static toObject(includeInstance: boolean, msg: Change): Change.AsObject;
    static serializeBinaryToWriter(message: Change, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Change;
    static deserializeBinaryFromReader(message: Change, reader: jspb.BinaryReader): Change;
  }

  export namespace Change {
    export type AsObject = {
      name: string,
      changeTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    }
  }

}

export class DescribeWasteRecordRequest extends jspb.Message {
  getName(): string;
  setName(value: string): DescribeWasteRecordRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DescribeWasteRecordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DescribeWasteRecordRequest): DescribeWasteRecordRequest.AsObject;
  static serializeBinaryToWriter(message: DescribeWasteRecordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DescribeWasteRecordRequest;
  static deserializeBinaryFromReader(message: DescribeWasteRecordRequest, reader: jspb.BinaryReader): DescribeWasteRecordRequest;
}

export namespace DescribeWasteRecordRequest {
  export type AsObject = {
    name: string,
  }
}

export class WasteRecordSupport extends jspb.Message {
  getResourceSupport(): types_info_pb.ResourceSupport | undefined;
  setResourceSupport(value?: types_info_pb.ResourceSupport): WasteRecordSupport;
  hasResourceSupport(): boolean;
  clearResourceSupport(): WasteRecordSupport;

  getUnit(): string;
  setUnit(value: string): WasteRecordSupport;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WasteRecordSupport.AsObject;
  static toObject(includeInstance: boolean, msg: WasteRecordSupport): WasteRecordSupport.AsObject;
  static serializeBinaryToWriter(message: WasteRecordSupport, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WasteRecordSupport;
  static deserializeBinaryFromReader(message: WasteRecordSupport, reader: jspb.BinaryReader): WasteRecordSupport;
}

export namespace WasteRecordSupport {
  export type AsObject = {
    resourceSupport?: types_info_pb.ResourceSupport.AsObject,
    unit: string,
  }
}

