import * as jspb from 'google-protobuf'

import * as google_protobuf_field_mask_pb from 'google-protobuf/google/protobuf/field_mask_pb'; // proto import: "google/protobuf/field_mask.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as actor_pb from './actor_pb'; // proto import: "actor.proto"
import * as types_info_pb from '@smart-core-os/sc-api-grpc-web/types/info_pb'; // proto import: "types/info.proto"


export class TransportState extends jspb.Message {
  getCurrentLocation(): string;
  setCurrentLocation(value: string): TransportState;

  getNextDestination(): string;
  setNextDestination(value: string): TransportState;

  getMovingDirection(): TransportState.Direction;
  setMovingDirection(value: TransportState.Direction): TransportState;

  getLoad(): number;
  setLoad(value: number): TransportState;

  getDoorStatusList(): Array<TransportState.DoorStatus>;
  setDoorStatusList(value: Array<TransportState.DoorStatus>): TransportState;
  clearDoorStatusList(): TransportState;
  addDoorStatus(value: TransportState.DoorStatus, index?: number): TransportState;

  getMode(): TransportState.Mode;
  setMode(value: TransportState.Mode): TransportState;

  getFaultsList(): Array<TransportState.Fault>;
  setFaultsList(value: Array<TransportState.Fault>): TransportState;
  clearFaultsList(): TransportState;
  addFaults(value: TransportState.Fault, index?: number): TransportState;

  getPassengerAlarm(): boolean;
  setPassengerAlarm(value: boolean): TransportState;

  getSpeed(): number;
  setSpeed(value: number): TransportState;

  getSupportedDestinationsList(): Array<number>;
  setSupportedDestinationsList(value: Array<number>): TransportState;
  clearSupportedDestinationsList(): TransportState;
  addSupportedDestinations(value: number, index?: number): TransportState;

  getActive(): boolean;
  setActive(value: boolean): TransportState;

  getActor(): actor_pb.Actor | undefined;
  setActor(value?: actor_pb.Actor): TransportState;
  hasActor(): boolean;
  clearActor(): TransportState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TransportState.AsObject;
  static toObject(includeInstance: boolean, msg: TransportState): TransportState.AsObject;
  static serializeBinaryToWriter(message: TransportState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TransportState;
  static deserializeBinaryFromReader(message: TransportState, reader: jspb.BinaryReader): TransportState;
}

export namespace TransportState {
  export type AsObject = {
    currentLocation: string,
    nextDestination: string,
    movingDirection: TransportState.Direction,
    load: number,
    doorStatusList: Array<TransportState.DoorStatus>,
    mode: TransportState.Mode,
    faultsList: Array<TransportState.Fault>,
    passengerAlarm: boolean,
    speed: number,
    supportedDestinationsList: Array<number>,
    active: boolean,
    actor?: actor_pb.Actor.AsObject,
  }

  export enum Direction { 
    DIRECTION_UNKNOWN = 0,
    NO_DIRECTION = 1,
    UP = 2,
    DOWN = 3,
    HORIZONTAL = 4,
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

export class TransportSupport extends jspb.Message {
  getResourceSupport(): types_info_pb.ResourceSupport | undefined;
  setResourceSupport(value?: types_info_pb.ResourceSupport): TransportSupport;
  hasResourceSupport(): boolean;
  clearResourceSupport(): TransportSupport;

  getLoadUnit(): string;
  setLoadUnit(value: string): TransportSupport;

  getDoorsList(): Array<TransportSupport.DoorInfo>;
  setDoorsList(value: Array<TransportSupport.DoorInfo>): TransportSupport;
  clearDoorsList(): TransportSupport;
  addDoors(value?: TransportSupport.DoorInfo, index?: number): TransportSupport.DoorInfo;

  getMaxLoad(): number;
  setMaxLoad(value: number): TransportSupport;

  getSpeedUnit(): string;
  setSpeedUnit(value: string): TransportSupport;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TransportSupport.AsObject;
  static toObject(includeInstance: boolean, msg: TransportSupport): TransportSupport.AsObject;
  static serializeBinaryToWriter(message: TransportSupport, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TransportSupport;
  static deserializeBinaryFromReader(message: TransportSupport, reader: jspb.BinaryReader): TransportSupport;
}

export namespace TransportSupport {
  export type AsObject = {
    resourceSupport?: types_info_pb.ResourceSupport.AsObject,
    loadUnit: string,
    doorsList: Array<TransportSupport.DoorInfo.AsObject>,
    maxLoad: number,
    speedUnit: string,
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

export class GetTransportStateRequest extends jspb.Message {
  getName(): string;
  setName(value: string): GetTransportStateRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): GetTransportStateRequest;
  hasReadMask(): boolean;
  clearReadMask(): GetTransportStateRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetTransportStateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetTransportStateRequest): GetTransportStateRequest.AsObject;
  static serializeBinaryToWriter(message: GetTransportStateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetTransportStateRequest;
  static deserializeBinaryFromReader(message: GetTransportStateRequest, reader: jspb.BinaryReader): GetTransportStateRequest;
}

export namespace GetTransportStateRequest {
  export type AsObject = {
    name: string,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
  }
}

export class DescribeTransportRequest extends jspb.Message {
  getName(): string;
  setName(value: string): DescribeTransportRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DescribeTransportRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DescribeTransportRequest): DescribeTransportRequest.AsObject;
  static serializeBinaryToWriter(message: DescribeTransportRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DescribeTransportRequest;
  static deserializeBinaryFromReader(message: DescribeTransportRequest, reader: jspb.BinaryReader): DescribeTransportRequest;
}

export namespace DescribeTransportRequest {
  export type AsObject = {
    name: string,
  }
}

export class PullTransportStateRequest extends jspb.Message {
  getName(): string;
  setName(value: string): PullTransportStateRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): PullTransportStateRequest;
  hasReadMask(): boolean;
  clearReadMask(): PullTransportStateRequest;

  getUpdatesOnly(): boolean;
  setUpdatesOnly(value: boolean): PullTransportStateRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullTransportStateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PullTransportStateRequest): PullTransportStateRequest.AsObject;
  static serializeBinaryToWriter(message: PullTransportStateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullTransportStateRequest;
  static deserializeBinaryFromReader(message: PullTransportStateRequest, reader: jspb.BinaryReader): PullTransportStateRequest;
}

export namespace PullTransportStateRequest {
  export type AsObject = {
    name: string,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
    updatesOnly: boolean,
  }
}

export class PullTransportStateResponse extends jspb.Message {
  getChangesList(): Array<PullTransportStateResponse.Change>;
  setChangesList(value: Array<PullTransportStateResponse.Change>): PullTransportStateResponse;
  clearChangesList(): PullTransportStateResponse;
  addChanges(value?: PullTransportStateResponse.Change, index?: number): PullTransportStateResponse.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullTransportStateResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PullTransportStateResponse): PullTransportStateResponse.AsObject;
  static serializeBinaryToWriter(message: PullTransportStateResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullTransportStateResponse;
  static deserializeBinaryFromReader(message: PullTransportStateResponse, reader: jspb.BinaryReader): PullTransportStateResponse;
}

export namespace PullTransportStateResponse {
  export type AsObject = {
    changesList: Array<PullTransportStateResponse.Change.AsObject>,
  }

  export class Change extends jspb.Message {
    getName(): string;
    setName(value: string): Change;

    getChangeTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setChangeTime(value?: google_protobuf_timestamp_pb.Timestamp): Change;
    hasChangeTime(): boolean;
    clearChangeTime(): Change;

    getTransportState(): TransportState | undefined;
    setTransportState(value?: TransportState): Change;
    hasTransportState(): boolean;
    clearTransportState(): Change;

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
      transportState?: TransportState.AsObject,
    }
  }

}

