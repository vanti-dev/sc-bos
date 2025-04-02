import * as jspb from 'google-protobuf'

import * as google_protobuf_field_mask_pb from 'google-protobuf/google/protobuf/field_mask_pb'; // proto import: "google/protobuf/field_mask.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"


export class AnprEvent extends jspb.Message {
  getEventTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setEventTime(value?: google_protobuf_timestamp_pb.Timestamp): AnprEvent;
  hasEventTime(): boolean;
  clearEventTime(): AnprEvent;

  getRegistrationPlate(): string;
  setRegistrationPlate(value: string): AnprEvent;

  getCountry(): string;
  setCountry(value: string): AnprEvent;

  getArea(): string;
  setArea(value: string): AnprEvent;
  hasArea(): boolean;
  clearArea(): AnprEvent;

  getConfidence(): number;
  setConfidence(value: number): AnprEvent;
  hasConfidence(): boolean;
  clearConfidence(): AnprEvent;

  getPlateType(): string;
  setPlateType(value: string): AnprEvent;
  hasPlateType(): boolean;
  clearPlateType(): AnprEvent;

  getYear(): string;
  setYear(value: string): AnprEvent;
  hasYear(): boolean;
  clearYear(): AnprEvent;

  getVehicleInfo(): AnprEvent.VehicleInfo | undefined;
  setVehicleInfo(value?: AnprEvent.VehicleInfo): AnprEvent;
  hasVehicleInfo(): boolean;
  clearVehicleInfo(): AnprEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AnprEvent.AsObject;
  static toObject(includeInstance: boolean, msg: AnprEvent): AnprEvent.AsObject;
  static serializeBinaryToWriter(message: AnprEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AnprEvent;
  static deserializeBinaryFromReader(message: AnprEvent, reader: jspb.BinaryReader): AnprEvent;
}

export namespace AnprEvent {
  export type AsObject = {
    eventTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    registrationPlate: string,
    country: string,
    area?: string,
    confidence?: number,
    plateType?: string,
    year?: string,
    vehicleInfo?: AnprEvent.VehicleInfo.AsObject,
  }

  export class VehicleInfo extends jspb.Message {
    getVehicleType(): string;
    setVehicleType(value: string): VehicleInfo;

    getColour(): string;
    setColour(value: string): VehicleInfo;

    getMake(): string;
    setMake(value: string): VehicleInfo;

    getModel(): string;
    setModel(value: string): VehicleInfo;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): VehicleInfo.AsObject;
    static toObject(includeInstance: boolean, msg: VehicleInfo): VehicleInfo.AsObject;
    static serializeBinaryToWriter(message: VehicleInfo, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): VehicleInfo;
    static deserializeBinaryFromReader(message: VehicleInfo, reader: jspb.BinaryReader): VehicleInfo;
  }

  export namespace VehicleInfo {
    export type AsObject = {
      vehicleType: string,
      colour: string,
      make: string,
      model: string,
    }
  }


  export enum AreaCase { 
    _AREA_NOT_SET = 0,
    AREA = 4,
  }

  export enum ConfidenceCase { 
    _CONFIDENCE_NOT_SET = 0,
    CONFIDENCE = 5,
  }

  export enum PlateTypeCase { 
    _PLATE_TYPE_NOT_SET = 0,
    PLATE_TYPE = 6,
  }

  export enum YearCase { 
    _YEAR_NOT_SET = 0,
    YEAR = 7,
  }

  export enum VehicleInfoCase { 
    _VEHICLE_INFO_NOT_SET = 0,
    VEHICLE_INFO = 8,
  }
}

export class GetLastEventRequest extends jspb.Message {
  getName(): string;
  setName(value: string): GetLastEventRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): GetLastEventRequest;
  hasReadMask(): boolean;
  clearReadMask(): GetLastEventRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLastEventRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetLastEventRequest): GetLastEventRequest.AsObject;
  static serializeBinaryToWriter(message: GetLastEventRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLastEventRequest;
  static deserializeBinaryFromReader(message: GetLastEventRequest, reader: jspb.BinaryReader): GetLastEventRequest;
}

export namespace GetLastEventRequest {
  export type AsObject = {
    name: string,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
  }
}

export class PullEventsRequest extends jspb.Message {
  getName(): string;
  setName(value: string): PullEventsRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): PullEventsRequest;
  hasReadMask(): boolean;
  clearReadMask(): PullEventsRequest;

  getUpdatesOnly(): boolean;
  setUpdatesOnly(value: boolean): PullEventsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullEventsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PullEventsRequest): PullEventsRequest.AsObject;
  static serializeBinaryToWriter(message: PullEventsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullEventsRequest;
  static deserializeBinaryFromReader(message: PullEventsRequest, reader: jspb.BinaryReader): PullEventsRequest;
}

export namespace PullEventsRequest {
  export type AsObject = {
    name: string,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
    updatesOnly: boolean,
  }
}

export class PullEventsResponse extends jspb.Message {
  getChangesList(): Array<PullEventsResponse.Change>;
  setChangesList(value: Array<PullEventsResponse.Change>): PullEventsResponse;
  clearChangesList(): PullEventsResponse;
  addChanges(value?: PullEventsResponse.Change, index?: number): PullEventsResponse.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullEventsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PullEventsResponse): PullEventsResponse.AsObject;
  static serializeBinaryToWriter(message: PullEventsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullEventsResponse;
  static deserializeBinaryFromReader(message: PullEventsResponse, reader: jspb.BinaryReader): PullEventsResponse;
}

export namespace PullEventsResponse {
  export type AsObject = {
    changesList: Array<PullEventsResponse.Change.AsObject>,
  }

  export class Change extends jspb.Message {
    getName(): string;
    setName(value: string): Change;

    getChangeTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setChangeTime(value?: google_protobuf_timestamp_pb.Timestamp): Change;
    hasChangeTime(): boolean;
    clearChangeTime(): Change;

    getAnprEvent(): AnprEvent | undefined;
    setAnprEvent(value?: AnprEvent): Change;
    hasAnprEvent(): boolean;
    clearAnprEvent(): Change;

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
      anprEvent?: AnprEvent.AsObject,
    }
  }

}

