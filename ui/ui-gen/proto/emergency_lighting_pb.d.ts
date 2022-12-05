import * as jspb from 'google-protobuf'

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb';


export class EmergencyLight extends jspb.Message {
  getName(): string;

  setName(value: string): EmergencyLight;

  getUpdateTime(): google_protobuf_timestamp_pb.Timestamp | undefined;

  setUpdateTime(value?: google_protobuf_timestamp_pb.Timestamp): EmergencyLight;

  hasUpdateTime(): boolean;

  clearUpdateTime(): EmergencyLight;

  getFaultsList(): Array<EmergencyLightFault>;

  setFaultsList(value: Array<EmergencyLightFault>): EmergencyLight;

  clearFaultsList(): EmergencyLight;

  addFaults(value: EmergencyLightFault, index?: number): EmergencyLight;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): EmergencyLight.AsObject;

  static toObject(includeInstance: boolean, msg: EmergencyLight): EmergencyLight.AsObject;

  static serializeBinaryToWriter(message: EmergencyLight, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): EmergencyLight;

  static deserializeBinaryFromReader(message: EmergencyLight, reader: jspb.BinaryReader): EmergencyLight;
}

export namespace EmergencyLight {
  export type AsObject = {
    name: string,
    updateTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    faultsList: Array<EmergencyLightFault>,
  }
}

export class EmergencyLightingEvent extends jspb.Message {
  getName(): string;

  setName(value: string): EmergencyLightingEvent;

  getId(): string;

  setId(value: string): EmergencyLightingEvent;

  getTimestamp(): google_protobuf_timestamp_pb.Timestamp | undefined;

  setTimestamp(value?: google_protobuf_timestamp_pb.Timestamp): EmergencyLightingEvent;

  hasTimestamp(): boolean;

  clearTimestamp(): EmergencyLightingEvent;

  getDurationTestPass(): EmergencyLightingEvent.DurationTestPass | undefined;

  setDurationTestPass(value?: EmergencyLightingEvent.DurationTestPass): EmergencyLightingEvent;

  hasDurationTestPass(): boolean;

  clearDurationTestPass(): EmergencyLightingEvent;

  getFunctionTestPass(): EmergencyLightingEvent.FunctionTestPass | undefined;

  setFunctionTestPass(value?: EmergencyLightingEvent.FunctionTestPass): EmergencyLightingEvent;

  hasFunctionTestPass(): boolean;

  clearFunctionTestPass(): EmergencyLightingEvent;

  getStatusReport(): EmergencyLightingEvent.StatusReport | undefined;

  setStatusReport(value?: EmergencyLightingEvent.StatusReport): EmergencyLightingEvent;

  hasStatusReport(): boolean;

  clearStatusReport(): EmergencyLightingEvent;

  getEventCase(): EmergencyLightingEvent.EventCase;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): EmergencyLightingEvent.AsObject;

  static toObject(includeInstance: boolean, msg: EmergencyLightingEvent): EmergencyLightingEvent.AsObject;

  static serializeBinaryToWriter(message: EmergencyLightingEvent, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): EmergencyLightingEvent;

  static deserializeBinaryFromReader(message: EmergencyLightingEvent, reader: jspb.BinaryReader): EmergencyLightingEvent;
}

export namespace EmergencyLightingEvent {
  export type AsObject = {
    name: string,
    id: string,
    timestamp?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    durationTestPass?: EmergencyLightingEvent.DurationTestPass.AsObject,
    functionTestPass?: EmergencyLightingEvent.FunctionTestPass.AsObject,
    statusReport?: EmergencyLightingEvent.StatusReport.AsObject,
  }

  export class DurationTestPass extends jspb.Message {
    getAchievedDuration(): google_protobuf_duration_pb.Duration | undefined;

    setAchievedDuration(value?: google_protobuf_duration_pb.Duration): DurationTestPass;

    hasAchievedDuration(): boolean;

    clearAchievedDuration(): DurationTestPass;

    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): DurationTestPass.AsObject;

    static toObject(includeInstance: boolean, msg: DurationTestPass): DurationTestPass.AsObject;

    static serializeBinaryToWriter(message: DurationTestPass, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): DurationTestPass;

    static deserializeBinaryFromReader(message: DurationTestPass, reader: jspb.BinaryReader): DurationTestPass;
  }

  export namespace DurationTestPass {
    export type AsObject = {
      achievedDuration?: google_protobuf_duration_pb.Duration.AsObject,
    }
  }


  export class FunctionTestPass extends jspb.Message {
    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): FunctionTestPass.AsObject;

    static toObject(includeInstance: boolean, msg: FunctionTestPass): FunctionTestPass.AsObject;

    static serializeBinaryToWriter(message: FunctionTestPass, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): FunctionTestPass;

    static deserializeBinaryFromReader(message: FunctionTestPass, reader: jspb.BinaryReader): FunctionTestPass;
  }

  export namespace FunctionTestPass {
    export type AsObject = {}
  }


  export class StatusReport extends jspb.Message {
    getFaultsList(): Array<EmergencyLightFault>;

    setFaultsList(value: Array<EmergencyLightFault>): StatusReport;

    clearFaultsList(): StatusReport;

    addFaults(value: EmergencyLightFault, index?: number): StatusReport;

    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): StatusReport.AsObject;

    static toObject(includeInstance: boolean, msg: StatusReport): StatusReport.AsObject;

    static serializeBinaryToWriter(message: StatusReport, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): StatusReport;

    static deserializeBinaryFromReader(message: StatusReport, reader: jspb.BinaryReader): StatusReport;
  }

  export namespace StatusReport {
    export type AsObject = {
      faultsList: Array<EmergencyLightFault>,
    }
  }


  export enum EventCase {
    EVENT_NOT_SET = 0,
    DURATION_TEST_PASS = 4,
    FUNCTION_TEST_PASS = 5,
    STATUS_REPORT = 6,
  }
}

export class GetEmergencyLightRequest extends jspb.Message {
  getName(): string;

  setName(value: string): GetEmergencyLightRequest;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): GetEmergencyLightRequest.AsObject;

  static toObject(includeInstance: boolean, msg: GetEmergencyLightRequest): GetEmergencyLightRequest.AsObject;

  static serializeBinaryToWriter(message: GetEmergencyLightRequest, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): GetEmergencyLightRequest;

  static deserializeBinaryFromReader(message: GetEmergencyLightRequest, reader: jspb.BinaryReader): GetEmergencyLightRequest;
}

export namespace GetEmergencyLightRequest {
  export type AsObject = {
    name: string,
  }
}

export class ListEmergencyLightsRequest extends jspb.Message {
  getPageSize(): number;

  setPageSize(value: number): ListEmergencyLightsRequest;

  getPageToken(): string;

  setPageToken(value: string): ListEmergencyLightsRequest;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): ListEmergencyLightsRequest.AsObject;

  static toObject(includeInstance: boolean, msg: ListEmergencyLightsRequest): ListEmergencyLightsRequest.AsObject;

  static serializeBinaryToWriter(message: ListEmergencyLightsRequest, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): ListEmergencyLightsRequest;

  static deserializeBinaryFromReader(message: ListEmergencyLightsRequest, reader: jspb.BinaryReader): ListEmergencyLightsRequest;
}

export namespace ListEmergencyLightsRequest {
  export type AsObject = {
    pageSize: number,
    pageToken: string,
  }
}

export class ListEmergencyLightsResponse extends jspb.Message {
  getEmergencyLightsList(): Array<EmergencyLight>;

  setEmergencyLightsList(value: Array<EmergencyLight>): ListEmergencyLightsResponse;

  clearEmergencyLightsList(): ListEmergencyLightsResponse;

  addEmergencyLights(value?: EmergencyLight, index?: number): EmergencyLight;

  getNextPageToken(): string;

  setNextPageToken(value: string): ListEmergencyLightsResponse;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): ListEmergencyLightsResponse.AsObject;

  static toObject(includeInstance: boolean, msg: ListEmergencyLightsResponse): ListEmergencyLightsResponse.AsObject;

  static serializeBinaryToWriter(message: ListEmergencyLightsResponse, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): ListEmergencyLightsResponse;

  static deserializeBinaryFromReader(message: ListEmergencyLightsResponse, reader: jspb.BinaryReader): ListEmergencyLightsResponse;
}

export namespace ListEmergencyLightsResponse {
  export type AsObject = {
    emergencyLightsList: Array<EmergencyLight.AsObject>,
    nextPageToken: string,
  }
}

export class ListEmergencyLightEventsRequest extends jspb.Message {
  getPageSize(): number;

  setPageSize(value: number): ListEmergencyLightEventsRequest;

  getPageToken(): string;

  setPageToken(value: string): ListEmergencyLightEventsRequest;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): ListEmergencyLightEventsRequest.AsObject;

  static toObject(includeInstance: boolean, msg: ListEmergencyLightEventsRequest): ListEmergencyLightEventsRequest.AsObject;

  static serializeBinaryToWriter(message: ListEmergencyLightEventsRequest, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): ListEmergencyLightEventsRequest;

  static deserializeBinaryFromReader(message: ListEmergencyLightEventsRequest, reader: jspb.BinaryReader): ListEmergencyLightEventsRequest;
}

export namespace ListEmergencyLightEventsRequest {
  export type AsObject = {
    pageSize: number,
    pageToken: string,
  }
}

export class ListEmergencyLightEventsResponse extends jspb.Message {
  getEventsList(): Array<EmergencyLightingEvent>;

  setEventsList(value: Array<EmergencyLightingEvent>): ListEmergencyLightEventsResponse;

  clearEventsList(): ListEmergencyLightEventsResponse;

  addEvents(value?: EmergencyLightingEvent, index?: number): EmergencyLightingEvent;

  getNextPageToken(): string;

  setNextPageToken(value: string): ListEmergencyLightEventsResponse;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): ListEmergencyLightEventsResponse.AsObject;

  static toObject(includeInstance: boolean, msg: ListEmergencyLightEventsResponse): ListEmergencyLightEventsResponse.AsObject;

  static serializeBinaryToWriter(message: ListEmergencyLightEventsResponse, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): ListEmergencyLightEventsResponse;

  static deserializeBinaryFromReader(message: ListEmergencyLightEventsResponse, reader: jspb.BinaryReader): ListEmergencyLightEventsResponse;
}

export namespace ListEmergencyLightEventsResponse {
  export type AsObject = {
    eventsList: Array<EmergencyLightingEvent.AsObject>,
    nextPageToken: string,
  }
}

export class GetReportCSVRequest extends jspb.Message {
  getIncludeHeader(): boolean;

  setIncludeHeader(value: boolean): GetReportCSVRequest;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): GetReportCSVRequest.AsObject;

  static toObject(includeInstance: boolean, msg: GetReportCSVRequest): GetReportCSVRequest.AsObject;

  static serializeBinaryToWriter(message: GetReportCSVRequest, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): GetReportCSVRequest;

  static deserializeBinaryFromReader(message: GetReportCSVRequest, reader: jspb.BinaryReader): GetReportCSVRequest;
}

export namespace GetReportCSVRequest {
  export type AsObject = {
    includeHeader: boolean,
  }
}

export class ReportCSV extends jspb.Message {
  getCsv(): Uint8Array | string;

  getCsv_asU8(): Uint8Array;

  getCsv_asB64(): string;

  setCsv(value: Uint8Array | string): ReportCSV;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): ReportCSV.AsObject;

  static toObject(includeInstance: boolean, msg: ReportCSV): ReportCSV.AsObject;

  static serializeBinaryToWriter(message: ReportCSV, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): ReportCSV;

  static deserializeBinaryFromReader(message: ReportCSV, reader: jspb.BinaryReader): ReportCSV;
}

export namespace ReportCSV {
  export type AsObject = {
    csv: Uint8Array | string,
  }
}

export enum EmergencyLightFault {
  FAULT_UNSPECIFIED = 0,
  DURATION_TEST_FAILED = 1,
  FUNCTION_TEST_FAILED = 2,
  BATTERY_FAULT = 3,
  LAMP_FAULT = 4,
  COMMUNICATION_FAILURE = 5,
  OTHER_FAULT = 6,
}
