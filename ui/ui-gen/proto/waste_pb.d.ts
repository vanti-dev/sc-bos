import * as jspb from 'google-protobuf'

import * as google_protobuf_field_mask_pb from 'google-protobuf/google/protobuf/field_mask_pb'; // proto import: "google/protobuf/field_mask.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"


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

export class GetWasteRecordsRequest extends jspb.Message {
  getName(): string;
  setName(value: string): GetWasteRecordsRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): GetWasteRecordsRequest;
  hasReadMask(): boolean;
  clearReadMask(): GetWasteRecordsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetWasteRecordsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetWasteRecordsRequest): GetWasteRecordsRequest.AsObject;
  static serializeBinaryToWriter(message: GetWasteRecordsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetWasteRecordsRequest;
  static deserializeBinaryFromReader(message: GetWasteRecordsRequest, reader: jspb.BinaryReader): GetWasteRecordsRequest;
}

export namespace GetWasteRecordsRequest {
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

