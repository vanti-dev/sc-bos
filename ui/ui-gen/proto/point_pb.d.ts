import * as jspb from 'google-protobuf'

import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb';
import * as google_protobuf_field_mask_pb from 'google-protobuf/google/protobuf/field_mask_pb';


export class GetPointsRequest extends jspb.Message {
  getName(): string;
  setName(value: string): GetPointsRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): GetPointsRequest;
  hasReadMask(): boolean;
  clearReadMask(): GetPointsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPointsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPointsRequest): GetPointsRequest.AsObject;
  static serializeBinaryToWriter(message: GetPointsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPointsRequest;
  static deserializeBinaryFromReader(message: GetPointsRequest, reader: jspb.BinaryReader): GetPointsRequest;
}

export namespace GetPointsRequest {
  export type AsObject = {
    name: string,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
  }
}

export class PullPointsRequest extends jspb.Message {
  getName(): string;
  setName(value: string): PullPointsRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): PullPointsRequest;
  hasReadMask(): boolean;
  clearReadMask(): PullPointsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullPointsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PullPointsRequest): PullPointsRequest.AsObject;
  static serializeBinaryToWriter(message: PullPointsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullPointsRequest;
  static deserializeBinaryFromReader(message: PullPointsRequest, reader: jspb.BinaryReader): PullPointsRequest;
}

export namespace PullPointsRequest {
  export type AsObject = {
    name: string,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
  }
}

export class PullPointsResponse extends jspb.Message {
  getPoints(): Points | undefined;
  setPoints(value?: Points): PullPointsResponse;
  hasPoints(): boolean;
  clearPoints(): PullPointsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullPointsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PullPointsResponse): PullPointsResponse.AsObject;
  static serializeBinaryToWriter(message: PullPointsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullPointsResponse;
  static deserializeBinaryFromReader(message: PullPointsResponse, reader: jspb.BinaryReader): PullPointsResponse;
}

export namespace PullPointsResponse {
  export type AsObject = {
    points?: Points.AsObject,
  }
}

export class Points extends jspb.Message {
  getValues(): google_protobuf_struct_pb.Struct | undefined;
  setValues(value?: google_protobuf_struct_pb.Struct): Points;
  hasValues(): boolean;
  clearValues(): Points;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Points.AsObject;
  static toObject(includeInstance: boolean, msg: Points): Points.AsObject;
  static serializeBinaryToWriter(message: Points, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Points;
  static deserializeBinaryFromReader(message: Points, reader: jspb.BinaryReader): Points;
}

export namespace Points {
  export type AsObject = {
    values?: google_protobuf_struct_pb.Struct.AsObject,
  }
}

export class DescribePointsRequest extends jspb.Message {
  getName(): string;
  setName(value: string): DescribePointsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DescribePointsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DescribePointsRequest): DescribePointsRequest.AsObject;
  static serializeBinaryToWriter(message: DescribePointsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DescribePointsRequest;
  static deserializeBinaryFromReader(message: DescribePointsRequest, reader: jspb.BinaryReader): DescribePointsRequest;
}

export namespace DescribePointsRequest {
  export type AsObject = {
    name: string,
  }
}

export class PointsSupport extends jspb.Message {
  getPointsList(): Array<PointMetadata>;
  setPointsList(value: Array<PointMetadata>): PointsSupport;
  clearPointsList(): PointsSupport;
  addPoints(value?: PointMetadata, index?: number): PointMetadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PointsSupport.AsObject;
  static toObject(includeInstance: boolean, msg: PointsSupport): PointsSupport.AsObject;
  static serializeBinaryToWriter(message: PointsSupport, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PointsSupport;
  static deserializeBinaryFromReader(message: PointsSupport, reader: jspb.BinaryReader): PointsSupport;
}

export namespace PointsSupport {
  export type AsObject = {
    pointsList: Array<PointMetadata.AsObject>,
  }
}

export class PointMetadata extends jspb.Message {
  getName(): string;
  setName(value: string): PointMetadata;

  getKind(): PointKind;
  setKind(value: PointKind): PointMetadata;

  getUnit(): PointUnit;
  setUnit(value: PointUnit): PointMetadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PointMetadata.AsObject;
  static toObject(includeInstance: boolean, msg: PointMetadata): PointMetadata.AsObject;
  static serializeBinaryToWriter(message: PointMetadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PointMetadata;
  static deserializeBinaryFromReader(message: PointMetadata, reader: jspb.BinaryReader): PointMetadata;
}

export namespace PointMetadata {
  export type AsObject = {
    name: string,
    kind: PointKind,
    unit: PointUnit,
  }
}

export enum PointKind { 
  POINT_KIND_UNSPECIFIED = 0,
  COUNT = 1,
  TEMPERATURE = 2,
}
export enum PointUnit { 
  POINT_UNIT_UNSPECIFIED = 0,
  CELSIUS = 1,
  KELVIN = 2,
  FAHRENHEIT = 3,
}
