import * as jspb from 'google-protobuf'

import * as google_protobuf_field_mask_pb from 'google-protobuf/google/protobuf/field_mask_pb'; // proto import: "google/protobuf/field_mask.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as types_info_pb from '@smart-core-os/sc-api-grpc-web/types/info_pb'; // proto import: "types/info.proto"


export class BearerState extends jspb.Message {
  getCurrentFloor(): number;
  setCurrentFloor(value: number): BearerState;

  getNextStoppingFloor(): number;
  setNextStoppingFloor(value: number): BearerState;

  getMovingDirection(): BearerState.Direction;
  setMovingDirection(value: BearerState.Direction): BearerState;

  getLoad(): number;
  setLoad(value: number): BearerState;

  getDoorStatusList(): Array<BearerState.DoorStatus>;
  setDoorStatusList(value: Array<BearerState.DoorStatus>): BearerState;
  clearDoorStatusList(): BearerState;
  addDoorStatus(value: BearerState.DoorStatus, index?: number): BearerState;

  getMode(): BearerState.Mode;
  setMode(value: BearerState.Mode): BearerState;

  getFaultsList(): Array<BearerState.Fault>;
  setFaultsList(value: Array<BearerState.Fault>): BearerState;
  clearFaultsList(): BearerState;
  addFaults(value: BearerState.Fault, index?: number): BearerState;

  getPassengerAlarm(): boolean;
  setPassengerAlarm(value: boolean): BearerState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BearerState.AsObject;
  static toObject(includeInstance: boolean, msg: BearerState): BearerState.AsObject;
  static serializeBinaryToWriter(message: BearerState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BearerState;
  static deserializeBinaryFromReader(message: BearerState, reader: jspb.BinaryReader): BearerState;
}

export namespace BearerState {
  export type AsObject = {
    currentFloor: number,
    nextStoppingFloor: number,
    movingDirection: BearerState.Direction,
    load: number,
    doorStatusList: Array<BearerState.DoorStatus>,
    mode: BearerState.Mode,
    faultsList: Array<BearerState.Fault>,
    passengerAlarm: boolean,
  }

  export enum Direction { 
    DIRECTION_UNKNOWN = 0,
    NO_DIRECTION = 1,
    UP = 2,
    DOWN = 3,
  }

  export enum DoorStatus { 
    DOOR_STATUS_UNKNOWN = 0,
    CLOSED = 1,
    OPEN = 2,
    OPENING = 3,
    CLOSING = 4,
    SAFETY_LOCKED = 5,
    LIMITED_OPENED = 6,
  }

  export enum Mode { 
    MODE_UNKNOWN = 0,
    NORMAL = 1,
    SERVICE_CONTROL = 2,
    FIREFIGHTER_CONTROL = 3,
    OUT_OF_SERVICE = 4,
    EMERGENCY_POWER = 5,
    VIP_CONTROL = 6,
    EARTHQUAKE_OPERATION = 7,
    FIRE_OPERATION = 8,
    ATTENDANT_CONTROL = 9,
    PARKING = 10,
    HOMING = 11,
    CABINET_RECALL = 12,
    OCCUPANT_EVACUATION = 13,
    FREIGHT = 14,
    FAILURE = 15,
    REDUCED_SPEED = 16,
    STORM_OPERATION = 17,
    HIGH_WIND_OPERATION = 18,
  }

  export enum Fault { 
    FAULT_UNKNOWN = 0,
    CONTROLLER_FAULT = 1,
    DRIVE_AND_MOTOR_FAULT = 2,
    MECHANICAL_COMPONENT_FAULT = 3,
    OVERSPEED_FAULT = 4,
    POWER_SUPPLY_FAULT = 5,
    SAFETY_DEVICE_FAULT = 6,
    CONTROLLER_SUPPLY_FAULT = 7,
    DRIVE_TEMPERATURE_EXCEEDED = 8,
    COMB_PLATE_FAULT = 9,
    GENERAL_FAULT = 10,
    DOOR_FAULT = 11,
    LEVELLING_FAULT = 12,
    SAFETY_CIRCUIT_BREAK_FAULT = 13,
    FAIL_TO_START = 14,
    ALARM_BUTTON = 15,
    DOOR_NOT_CLOSING = 16,
  }
}

export class BearerStateSupport extends jspb.Message {
  getResourceSupport(): types_info_pb.ResourceSupport | undefined;
  setResourceSupport(value?: types_info_pb.ResourceSupport): BearerStateSupport;
  hasResourceSupport(): boolean;
  clearResourceSupport(): BearerStateSupport;

  getUnit(): string;
  setUnit(value: string): BearerStateSupport;

  getDoorsList(): Array<BearerStateSupport.DoorInfo>;
  setDoorsList(value: Array<BearerStateSupport.DoorInfo>): BearerStateSupport;
  clearDoorsList(): BearerStateSupport;
  addDoors(value?: BearerStateSupport.DoorInfo, index?: number): BearerStateSupport.DoorInfo;

  getMaxLoad(): number;
  setMaxLoad(value: number): BearerStateSupport;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BearerStateSupport.AsObject;
  static toObject(includeInstance: boolean, msg: BearerStateSupport): BearerStateSupport.AsObject;
  static serializeBinaryToWriter(message: BearerStateSupport, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BearerStateSupport;
  static deserializeBinaryFromReader(message: BearerStateSupport, reader: jspb.BinaryReader): BearerStateSupport;
}

export namespace BearerStateSupport {
  export type AsObject = {
    resourceSupport?: types_info_pb.ResourceSupport.AsObject,
    unit: string,
    doorsList: Array<BearerStateSupport.DoorInfo.AsObject>,
    maxLoad: number,
  }

  export class DoorInfo extends jspb.Message {
    getId(): number;
    setId(value: number): DoorInfo;

    getDescription(): string;
    setDescription(value: string): DoorInfo;

    getDeck(): number;
    setDeck(value: number): DoorInfo;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DoorInfo.AsObject;
    static toObject(includeInstance: boolean, msg: DoorInfo): DoorInfo.AsObject;
    static serializeBinaryToWriter(message: DoorInfo, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DoorInfo;
    static deserializeBinaryFromReader(message: DoorInfo, reader: jspb.BinaryReader): DoorInfo;
  }

  export namespace DoorInfo {
    export type AsObject = {
      id: number,
      description: string,
      deck: number,
    }
  }

}

export class GetBearerStateRequest extends jspb.Message {
  getName(): string;
  setName(value: string): GetBearerStateRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): GetBearerStateRequest;
  hasReadMask(): boolean;
  clearReadMask(): GetBearerStateRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetBearerStateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetBearerStateRequest): GetBearerStateRequest.AsObject;
  static serializeBinaryToWriter(message: GetBearerStateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetBearerStateRequest;
  static deserializeBinaryFromReader(message: GetBearerStateRequest, reader: jspb.BinaryReader): GetBearerStateRequest;
}

export namespace GetBearerStateRequest {
  export type AsObject = {
    name: string,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
  }
}

export class DescribeBearerRequest extends jspb.Message {
  getName(): string;
  setName(value: string): DescribeBearerRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DescribeBearerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DescribeBearerRequest): DescribeBearerRequest.AsObject;
  static serializeBinaryToWriter(message: DescribeBearerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DescribeBearerRequest;
  static deserializeBinaryFromReader(message: DescribeBearerRequest, reader: jspb.BinaryReader): DescribeBearerRequest;
}

export namespace DescribeBearerRequest {
  export type AsObject = {
    name: string,
  }
}

export class PullBearerStateRequest extends jspb.Message {
  getName(): string;
  setName(value: string): PullBearerStateRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): PullBearerStateRequest;
  hasReadMask(): boolean;
  clearReadMask(): PullBearerStateRequest;

  getUpdatesOnly(): boolean;
  setUpdatesOnly(value: boolean): PullBearerStateRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullBearerStateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PullBearerStateRequest): PullBearerStateRequest.AsObject;
  static serializeBinaryToWriter(message: PullBearerStateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullBearerStateRequest;
  static deserializeBinaryFromReader(message: PullBearerStateRequest, reader: jspb.BinaryReader): PullBearerStateRequest;
}

export namespace PullBearerStateRequest {
  export type AsObject = {
    name: string,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
    updatesOnly: boolean,
  }
}

export class PullBearerStateResponse extends jspb.Message {
  getChangesList(): Array<PullBearerStateResponse.Change>;
  setChangesList(value: Array<PullBearerStateResponse.Change>): PullBearerStateResponse;
  clearChangesList(): PullBearerStateResponse;
  addChanges(value?: PullBearerStateResponse.Change, index?: number): PullBearerStateResponse.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullBearerStateResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PullBearerStateResponse): PullBearerStateResponse.AsObject;
  static serializeBinaryToWriter(message: PullBearerStateResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullBearerStateResponse;
  static deserializeBinaryFromReader(message: PullBearerStateResponse, reader: jspb.BinaryReader): PullBearerStateResponse;
}

export namespace PullBearerStateResponse {
  export type AsObject = {
    changesList: Array<PullBearerStateResponse.Change.AsObject>,
  }

  export class Change extends jspb.Message {
    getName(): string;
    setName(value: string): Change;

    getChangeTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setChangeTime(value?: google_protobuf_timestamp_pb.Timestamp): Change;
    hasChangeTime(): boolean;
    clearChangeTime(): Change;

    getBearerState(): BearerState | undefined;
    setBearerState(value?: BearerState): Change;
    hasBearerState(): boolean;
    clearBearerState(): Change;

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
      bearerState?: BearerState.AsObject,
    }
  }

}

