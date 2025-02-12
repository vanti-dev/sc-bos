import * as jspb from 'google-protobuf'

import * as google_protobuf_field_mask_pb from 'google-protobuf/google/protobuf/field_mask_pb'; // proto import: "google/protobuf/field_mask.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as types_info_pb from '@smart-core-os/sc-api-grpc-web/types/info_pb'; // proto import: "types/info.proto"


export class Transport extends jspb.Message {
  getCurrentLocation(): Transport.Location | undefined;
  setCurrentLocation(value?: Transport.Location): Transport;
  hasCurrentLocation(): boolean;
  clearCurrentLocation(): Transport;

  getNextDestinationsList(): Array<Transport.Location>;
  setNextDestinationsList(value: Array<Transport.Location>): Transport;
  clearNextDestinationsList(): Transport;
  addNextDestinations(value?: Transport.Location, index?: number): Transport.Location;

  getMovingDirection(): Transport.Direction;
  setMovingDirection(value: Transport.Direction): Transport;

  getLoad(): number;
  setLoad(value: number): Transport;

  getDoorStatusList(): Array<Transport.DoorStatus>;
  setDoorStatusList(value: Array<Transport.DoorStatus>): Transport;
  clearDoorStatusList(): Transport;
  addDoorStatus(value: Transport.DoorStatus, index?: number): Transport;

  getMode(): Transport.Mode;
  setMode(value: Transport.Mode): Transport;

  getFaultsList(): Array<Transport.Fault>;
  setFaultsList(value: Array<Transport.Fault>): Transport;
  clearFaultsList(): Transport;
  addFaults(value: Transport.Fault, index?: number): Transport;

  getPassengerAlarm(): boolean;
  setPassengerAlarm(value: boolean): Transport;

  getSpeed(): number;
  setSpeed(value: number): Transport;

  getSupportedDestinationsList(): Array<Transport.Location>;
  setSupportedDestinationsList(value: Array<Transport.Location>): Transport;
  clearSupportedDestinationsList(): Transport;
  addSupportedDestinations(value?: Transport.Location, index?: number): Transport.Location;

  getActive(): boolean;
  setActive(value: boolean): Transport;

  getPayloadsList(): Array<Transport.Payload>;
  setPayloadsList(value: Array<Transport.Payload>): Transport;
  clearPayloadsList(): Transport;
  addPayloads(value?: Transport.Payload, index?: number): Transport.Payload;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Transport.AsObject;
  static toObject(includeInstance: boolean, msg: Transport): Transport.AsObject;
  static serializeBinaryToWriter(message: Transport, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Transport;
  static deserializeBinaryFromReader(message: Transport, reader: jspb.BinaryReader): Transport;
}

export namespace Transport {
  export type AsObject = {
    currentLocation?: Transport.Location.AsObject,
    nextDestinationsList: Array<Transport.Location.AsObject>,
    movingDirection: Transport.Direction,
    load: number,
    doorStatusList: Array<Transport.DoorStatus>,
    mode: Transport.Mode,
    faultsList: Array<Transport.Fault>,
    passengerAlarm: boolean,
    speed: number,
    supportedDestinationsList: Array<Transport.Location.AsObject>,
    active: boolean,
    payloadsList: Array<Transport.Payload.AsObject>,
  }

  export class Location extends jspb.Message {
    getId(): string;
    setId(value: string): Location;

    getName(): string;
    setName(value: string): Location;

    getDescription(): string;
    setDescription(value: string): Location;

    getFloor(): number;
    setFloor(value: number): Location;

    getAttributesList(): Array<string>;
    setAttributesList(value: Array<string>): Location;
    clearAttributesList(): Location;
    addAttributes(value: string, index?: number): Location;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Location.AsObject;
    static toObject(includeInstance: boolean, msg: Location): Location.AsObject;
    static serializeBinaryToWriter(message: Location, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Location;
    static deserializeBinaryFromReader(message: Location, reader: jspb.BinaryReader): Location;
  }

  export namespace Location {
    export type AsObject = {
      id: string,
      name: string,
      description: string,
      floor: number,
      attributesList: Array<string>,
    }
  }


  export class Journey extends jspb.Message {
    getStart(): Transport.Location | undefined;
    setStart(value?: Transport.Location): Journey;
    hasStart(): boolean;
    clearStart(): Journey;

    getDestinationsList(): Array<Transport.Location>;
    setDestinationsList(value: Array<Transport.Location>): Journey;
    clearDestinationsList(): Journey;
    addDestinations(value?: Transport.Location, index?: number): Transport.Location;

    getReason(): string;
    setReason(value: string): Journey;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Journey.AsObject;
    static toObject(includeInstance: boolean, msg: Journey): Journey.AsObject;
    static serializeBinaryToWriter(message: Journey, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Journey;
    static deserializeBinaryFromReader(message: Journey, reader: jspb.BinaryReader): Journey;
  }

  export namespace Journey {
    export type AsObject = {
      start?: Transport.Location.AsObject,
      destinationsList: Array<Transport.Location.AsObject>,
      reason: string,
    }
  }


  export class Payload extends jspb.Message {
    getId(): string;
    setId(value: string): Payload;

    getDescription(): string;
    setDescription(value: string): Payload;

    getIntendedJourney(): Transport.Journey | undefined;
    setIntendedJourney(value?: Transport.Journey): Payload;
    hasIntendedJourney(): boolean;
    clearIntendedJourney(): Payload;

    getActualJourney(): Transport.Journey | undefined;
    setActualJourney(value?: Transport.Journey): Payload;
    hasActualJourney(): boolean;
    clearActualJourney(): Payload;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Payload.AsObject;
    static toObject(includeInstance: boolean, msg: Payload): Payload.AsObject;
    static serializeBinaryToWriter(message: Payload, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Payload;
    static deserializeBinaryFromReader(message: Payload, reader: jspb.BinaryReader): Payload;
  }

  export namespace Payload {
    export type AsObject = {
      id: string,
      description: string,
      intendedJourney?: Transport.Journey.AsObject,
      actualJourney?: Transport.Journey.AsObject,
    }
  }


  export enum Direction { 
    DIRECTION_UNSPECIFIED = 0,
    NO_DIRECTION = 1,
    UP = 2,
    DOWN = 3,
    IN = 4,
    OUT = 5,
    CLOCKWISE = 6,
    ANTICLOCKWISE = 7,
    FORWARD = 8,
    BACKWARD = 9,
    EAST = 10,
    WEST = 11,
    NORTH = 12,
    SOUTH = 13,
    LEFT = 14,
    RIGHT = 15,
    SIDEWAYS = 16,
  }

  export enum DoorStatus { 
    DOOR_STATUS_UNSPECIFIED = 0,
    CLOSED = 1,
    OPEN = 2,
    OPENING = 3,
    CLOSING = 4,
    SAFETY_LOCKED = 5,
    LIMITED_OPENED = 6,
  }

  export enum Mode { 
    MODE_UNSPECIFIED = 0,
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
    FAULT_UNSPECIFIED = 0,
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

export class GetTransportRequest extends jspb.Message {
  getName(): string;
  setName(value: string): GetTransportRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): GetTransportRequest;
  hasReadMask(): boolean;
  clearReadMask(): GetTransportRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetTransportRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetTransportRequest): GetTransportRequest.AsObject;
  static serializeBinaryToWriter(message: GetTransportRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetTransportRequest;
  static deserializeBinaryFromReader(message: GetTransportRequest, reader: jspb.BinaryReader): GetTransportRequest;
}

export namespace GetTransportRequest {
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

export class PullTransportRequest extends jspb.Message {
  getName(): string;
  setName(value: string): PullTransportRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): PullTransportRequest;
  hasReadMask(): boolean;
  clearReadMask(): PullTransportRequest;

  getUpdatesOnly(): boolean;
  setUpdatesOnly(value: boolean): PullTransportRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullTransportRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PullTransportRequest): PullTransportRequest.AsObject;
  static serializeBinaryToWriter(message: PullTransportRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullTransportRequest;
  static deserializeBinaryFromReader(message: PullTransportRequest, reader: jspb.BinaryReader): PullTransportRequest;
}

export namespace PullTransportRequest {
  export type AsObject = {
    name: string,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
    updatesOnly: boolean,
  }
}

export class PullTransportResponse extends jspb.Message {
  getChangesList(): Array<PullTransportResponse.Change>;
  setChangesList(value: Array<PullTransportResponse.Change>): PullTransportResponse;
  clearChangesList(): PullTransportResponse;
  addChanges(value?: PullTransportResponse.Change, index?: number): PullTransportResponse.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullTransportResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PullTransportResponse): PullTransportResponse.AsObject;
  static serializeBinaryToWriter(message: PullTransportResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullTransportResponse;
  static deserializeBinaryFromReader(message: PullTransportResponse, reader: jspb.BinaryReader): PullTransportResponse;
}

export namespace PullTransportResponse {
  export type AsObject = {
    changesList: Array<PullTransportResponse.Change.AsObject>,
  }

  export class Change extends jspb.Message {
    getName(): string;
    setName(value: string): Change;

    getChangeTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setChangeTime(value?: google_protobuf_timestamp_pb.Timestamp): Change;
    hasChangeTime(): boolean;
    clearChangeTime(): Change;

    getTransportState(): Transport | undefined;
    setTransportState(value?: Transport): Change;
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
      transportState?: Transport.AsObject,
    }
  }

}

