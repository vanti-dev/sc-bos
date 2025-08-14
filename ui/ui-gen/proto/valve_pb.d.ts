import * as jspb from 'google-protobuf'

import * as google_protobuf_field_mask_pb from 'google-protobuf/google/protobuf/field_mask_pb'; // proto import: "google/protobuf/field_mask.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"


export class Valve extends jspb.Message {
  getFlowRateSetPoint(): number;
  setFlowRateSetPoint(value: number): Valve;

  getFlowRateSensor(): number;
  setFlowRateSensor(value: number): Valve;

  getPumpageSensor(): number;
  setPumpageSensor(value: number): Valve;

  getPumpageSetPoint(): number;
  setPumpageSetPoint(value: number): Valve;

  getTemperatureSensor(): number;
  setTemperatureSensor(value: number): Valve;

  getTemperatureSetPoint(): number;
  setTemperatureSetPoint(value: number): Valve;

  getPressureSensor(): number;
  setPressureSensor(value: number): Valve;

  getPressureSetPoint(): number;
  setPressureSetPoint(value: number): Valve;

  getDriveFrequencySensor(): number;
  setDriveFrequencySensor(value: number): Valve;

  getDriveFrequencySetPoint(): number;
  setDriveFrequencySetPoint(value: number): Valve;

  getMode(): Valve.ValveMode;
  setMode(value: Valve.ValveMode): Valve;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Valve.AsObject;
  static toObject(includeInstance: boolean, msg: Valve): Valve.AsObject;
  static serializeBinaryToWriter(message: Valve, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Valve;
  static deserializeBinaryFromReader(message: Valve, reader: jspb.BinaryReader): Valve;
}

export namespace Valve {
  export type AsObject = {
    flowRateSetPoint: number,
    flowRateSensor: number,
    pumpageSensor: number,
    pumpageSetPoint: number,
    temperatureSensor: number,
    temperatureSetPoint: number,
    pressureSensor: number,
    pressureSetPoint: number,
    driveFrequencySensor: number,
    driveFrequencySetPoint: number,
    mode: Valve.ValveMode,
  }

  export enum ValveMode { 
    VALVE_MODE_UNKNOWN = 0,
    VALVE_MODE_FLOW = 1,
    VALVE_MODE_RETURN = 2,
    VALVE_MODE_BLOCKED = 3,
  }
}

export class GetValveRequest extends jspb.Message {
  getName(): string;
  setName(value: string): GetValveRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): GetValveRequest;
  hasReadMask(): boolean;
  clearReadMask(): GetValveRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetValveRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetValveRequest): GetValveRequest.AsObject;
  static serializeBinaryToWriter(message: GetValveRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetValveRequest;
  static deserializeBinaryFromReader(message: GetValveRequest, reader: jspb.BinaryReader): GetValveRequest;
}

export namespace GetValveRequest {
  export type AsObject = {
    name: string,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
  }
}

export class PullValveRequest extends jspb.Message {
  getName(): string;
  setName(value: string): PullValveRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): PullValveRequest;
  hasReadMask(): boolean;
  clearReadMask(): PullValveRequest;

  getUpdatesOnly(): boolean;
  setUpdatesOnly(value: boolean): PullValveRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullValveRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PullValveRequest): PullValveRequest.AsObject;
  static serializeBinaryToWriter(message: PullValveRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullValveRequest;
  static deserializeBinaryFromReader(message: PullValveRequest, reader: jspb.BinaryReader): PullValveRequest;
}

export namespace PullValveRequest {
  export type AsObject = {
    name: string,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
    updatesOnly: boolean,
  }
}

export class PullValveResponse extends jspb.Message {
  getChangesList(): Array<PullValveResponse.Change>;
  setChangesList(value: Array<PullValveResponse.Change>): PullValveResponse;
  clearChangesList(): PullValveResponse;
  addChanges(value?: PullValveResponse.Change, index?: number): PullValveResponse.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullValveResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PullValveResponse): PullValveResponse.AsObject;
  static serializeBinaryToWriter(message: PullValveResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullValveResponse;
  static deserializeBinaryFromReader(message: PullValveResponse, reader: jspb.BinaryReader): PullValveResponse;
}

export namespace PullValveResponse {
  export type AsObject = {
    changesList: Array<PullValveResponse.Change.AsObject>,
  }

  export class Change extends jspb.Message {
    getName(): string;
    setName(value: string): Change;

    getChangeTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setChangeTime(value?: google_protobuf_timestamp_pb.Timestamp): Change;
    hasChangeTime(): boolean;
    clearChangeTime(): Change;

    getValve(): Valve | undefined;
    setValve(value?: Valve): Change;
    hasValve(): boolean;
    clearValve(): Change;

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
      valve?: Valve.AsObject,
    }
  }

}

export class UpdateValveRequest extends jspb.Message {
  getName(): string;
  setName(value: string): UpdateValveRequest;

  getValve(): Valve | undefined;
  setValve(value?: Valve): UpdateValveRequest;
  hasValve(): boolean;
  clearValve(): UpdateValveRequest;

  getUpdateMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setUpdateMask(value?: google_protobuf_field_mask_pb.FieldMask): UpdateValveRequest;
  hasUpdateMask(): boolean;
  clearUpdateMask(): UpdateValveRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateValveRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateValveRequest): UpdateValveRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateValveRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateValveRequest;
  static deserializeBinaryFromReader(message: UpdateValveRequest, reader: jspb.BinaryReader): UpdateValveRequest;
}

export namespace UpdateValveRequest {
  export type AsObject = {
    name: string,
    valve?: Valve.AsObject,
    updateMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
  }
}

export class UpdateValveResponse extends jspb.Message {
  getValve(): Valve | undefined;
  setValve(value?: Valve): UpdateValveResponse;
  hasValve(): boolean;
  clearValve(): UpdateValveResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateValveResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateValveResponse): UpdateValveResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateValveResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateValveResponse;
  static deserializeBinaryFromReader(message: UpdateValveResponse, reader: jspb.BinaryReader): UpdateValveResponse;
}

export namespace UpdateValveResponse {
  export type AsObject = {
    valve?: Valve.AsObject,
  }
}

