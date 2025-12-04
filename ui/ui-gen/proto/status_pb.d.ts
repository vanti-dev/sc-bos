import * as jspb from 'google-protobuf'

import * as google_protobuf_field_mask_pb from 'google-protobuf/google/protobuf/field_mask_pb'; // proto import: "google/protobuf/field_mask.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as types_time_period_pb from '@smart-core-os/sc-api-grpc-web/types/time/period_pb'; // proto import: "types/time/period.proto"


export class StatusLog extends jspb.Message {
  getLevel(): StatusLog.Level;
  setLevel(value: StatusLog.Level): StatusLog;

  getDescription(): string;
  setDescription(value: string): StatusLog;

  getRecordTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setRecordTime(value?: google_protobuf_timestamp_pb.Timestamp): StatusLog;
  hasRecordTime(): boolean;
  clearRecordTime(): StatusLog;

  getProblemsList(): Array<StatusLog.Problem>;
  setProblemsList(value: Array<StatusLog.Problem>): StatusLog;
  clearProblemsList(): StatusLog;
  addProblems(value?: StatusLog.Problem, index?: number): StatusLog.Problem;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StatusLog.AsObject;
  static toObject(includeInstance: boolean, msg: StatusLog): StatusLog.AsObject;
  static serializeBinaryToWriter(message: StatusLog, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StatusLog;
  static deserializeBinaryFromReader(message: StatusLog, reader: jspb.BinaryReader): StatusLog;
}

export namespace StatusLog {
  export type AsObject = {
    level: StatusLog.Level;
    description: string;
    recordTime?: google_protobuf_timestamp_pb.Timestamp.AsObject;
    problemsList: Array<StatusLog.Problem.AsObject>;
  };

  export class Problem extends jspb.Message {
    getLevel(): StatusLog.Level;
    setLevel(value: StatusLog.Level): Problem;

    getDescription(): string;
    setDescription(value: string): Problem;

    getRecordTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setRecordTime(value?: google_protobuf_timestamp_pb.Timestamp): Problem;
    hasRecordTime(): boolean;
    clearRecordTime(): Problem;

    getName(): string;
    setName(value: string): Problem;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Problem.AsObject;
    static toObject(includeInstance: boolean, msg: Problem): Problem.AsObject;
    static serializeBinaryToWriter(message: Problem, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Problem;
    static deserializeBinaryFromReader(message: Problem, reader: jspb.BinaryReader): Problem;
  }

  export namespace Problem {
    export type AsObject = {
      level: StatusLog.Level;
      description: string;
      recordTime?: google_protobuf_timestamp_pb.Timestamp.AsObject;
      name: string;
    };
  }


  export enum Level {
    LEVEL_UNDEFINED = 0,
    NOMINAL = 1,
    NOTICE = 2,
    REDUCED_FUNCTION = 3,
    NON_FUNCTIONAL = 4,
    OFFLINE = 127,
  }
}

export class GetCurrentStatusRequest extends jspb.Message {
  getName(): string;
  setName(value: string): GetCurrentStatusRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): GetCurrentStatusRequest;
  hasReadMask(): boolean;
  clearReadMask(): GetCurrentStatusRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCurrentStatusRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCurrentStatusRequest): GetCurrentStatusRequest.AsObject;
  static serializeBinaryToWriter(message: GetCurrentStatusRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCurrentStatusRequest;
  static deserializeBinaryFromReader(message: GetCurrentStatusRequest, reader: jspb.BinaryReader): GetCurrentStatusRequest;
}

export namespace GetCurrentStatusRequest {
  export type AsObject = {
    name: string;
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject;
  };
}

export class PullCurrentStatusRequest extends jspb.Message {
  getName(): string;
  setName(value: string): PullCurrentStatusRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): PullCurrentStatusRequest;
  hasReadMask(): boolean;
  clearReadMask(): PullCurrentStatusRequest;

  getUpdatesOnly(): boolean;
  setUpdatesOnly(value: boolean): PullCurrentStatusRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullCurrentStatusRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PullCurrentStatusRequest): PullCurrentStatusRequest.AsObject;
  static serializeBinaryToWriter(message: PullCurrentStatusRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullCurrentStatusRequest;
  static deserializeBinaryFromReader(message: PullCurrentStatusRequest, reader: jspb.BinaryReader): PullCurrentStatusRequest;
}

export namespace PullCurrentStatusRequest {
  export type AsObject = {
    name: string;
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject;
    updatesOnly: boolean;
  };
}

export class PullCurrentStatusResponse extends jspb.Message {
  getChangesList(): Array<PullCurrentStatusResponse.Change>;
  setChangesList(value: Array<PullCurrentStatusResponse.Change>): PullCurrentStatusResponse;
  clearChangesList(): PullCurrentStatusResponse;
  addChanges(value?: PullCurrentStatusResponse.Change, index?: number): PullCurrentStatusResponse.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullCurrentStatusResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PullCurrentStatusResponse): PullCurrentStatusResponse.AsObject;
  static serializeBinaryToWriter(message: PullCurrentStatusResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullCurrentStatusResponse;
  static deserializeBinaryFromReader(message: PullCurrentStatusResponse, reader: jspb.BinaryReader): PullCurrentStatusResponse;
}

export namespace PullCurrentStatusResponse {
  export type AsObject = {
    changesList: Array<PullCurrentStatusResponse.Change.AsObject>;
  };

  export class Change extends jspb.Message {
    getName(): string;
    setName(value: string): Change;

    getChangeTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setChangeTime(value?: google_protobuf_timestamp_pb.Timestamp): Change;
    hasChangeTime(): boolean;
    clearChangeTime(): Change;

    getCurrentStatus(): StatusLog | undefined;
    setCurrentStatus(value?: StatusLog): Change;
    hasCurrentStatus(): boolean;
    clearCurrentStatus(): Change;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Change.AsObject;
    static toObject(includeInstance: boolean, msg: Change): Change.AsObject;
    static serializeBinaryToWriter(message: Change, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Change;
    static deserializeBinaryFromReader(message: Change, reader: jspb.BinaryReader): Change;
  }

  export namespace Change {
    export type AsObject = {
      name: string;
      changeTime?: google_protobuf_timestamp_pb.Timestamp.AsObject;
      currentStatus?: StatusLog.AsObject;
    };
  }

}

export class StatusLogRecord extends jspb.Message {
  getCurrentStatus(): StatusLog | undefined;
  setCurrentStatus(value?: StatusLog): StatusLogRecord;
  hasCurrentStatus(): boolean;
  clearCurrentStatus(): StatusLogRecord;

  getRecordTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setRecordTime(value?: google_protobuf_timestamp_pb.Timestamp): StatusLogRecord;
  hasRecordTime(): boolean;
  clearRecordTime(): StatusLogRecord;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StatusLogRecord.AsObject;
  static toObject(includeInstance: boolean, msg: StatusLogRecord): StatusLogRecord.AsObject;
  static serializeBinaryToWriter(message: StatusLogRecord, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StatusLogRecord;
  static deserializeBinaryFromReader(message: StatusLogRecord, reader: jspb.BinaryReader): StatusLogRecord;
}

export namespace StatusLogRecord {
  export type AsObject = {
    currentStatus?: StatusLog.AsObject;
    recordTime?: google_protobuf_timestamp_pb.Timestamp.AsObject;
  };
}

export class ListCurrentStatusHistoryRequest extends jspb.Message {
  getName(): string;
  setName(value: string): ListCurrentStatusHistoryRequest;

  getPeriod(): types_time_period_pb.Period | undefined;
  setPeriod(value?: types_time_period_pb.Period): ListCurrentStatusHistoryRequest;
  hasPeriod(): boolean;
  clearPeriod(): ListCurrentStatusHistoryRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): ListCurrentStatusHistoryRequest;
  hasReadMask(): boolean;
  clearReadMask(): ListCurrentStatusHistoryRequest;

  getPageSize(): number;
  setPageSize(value: number): ListCurrentStatusHistoryRequest;

  getPageToken(): string;
  setPageToken(value: string): ListCurrentStatusHistoryRequest;

  getOrderBy(): string;
  setOrderBy(value: string): ListCurrentStatusHistoryRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListCurrentStatusHistoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListCurrentStatusHistoryRequest): ListCurrentStatusHistoryRequest.AsObject;
  static serializeBinaryToWriter(message: ListCurrentStatusHistoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListCurrentStatusHistoryRequest;
  static deserializeBinaryFromReader(message: ListCurrentStatusHistoryRequest, reader: jspb.BinaryReader): ListCurrentStatusHistoryRequest;
}

export namespace ListCurrentStatusHistoryRequest {
  export type AsObject = {
    name: string;
    period?: types_time_period_pb.Period.AsObject;
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject;
    pageSize: number;
    pageToken: string;
    orderBy: string;
  };
}

export class ListCurrentStatusHistoryResponse extends jspb.Message {
  getCurrentStatusRecordsList(): Array<StatusLogRecord>;
  setCurrentStatusRecordsList(value: Array<StatusLogRecord>): ListCurrentStatusHistoryResponse;
  clearCurrentStatusRecordsList(): ListCurrentStatusHistoryResponse;
  addCurrentStatusRecords(value?: StatusLogRecord, index?: number): StatusLogRecord;

  getNextPageToken(): string;
  setNextPageToken(value: string): ListCurrentStatusHistoryResponse;

  getTotalSize(): number;
  setTotalSize(value: number): ListCurrentStatusHistoryResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListCurrentStatusHistoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListCurrentStatusHistoryResponse): ListCurrentStatusHistoryResponse.AsObject;
  static serializeBinaryToWriter(message: ListCurrentStatusHistoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListCurrentStatusHistoryResponse;
  static deserializeBinaryFromReader(message: ListCurrentStatusHistoryResponse, reader: jspb.BinaryReader): ListCurrentStatusHistoryResponse;
}

export namespace ListCurrentStatusHistoryResponse {
  export type AsObject = {
    currentStatusRecordsList: Array<StatusLogRecord.AsObject>;
    nextPageToken: string;
    totalSize: number;
  };
}

