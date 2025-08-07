import * as jspb from 'google-protobuf'

import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"
import * as google_protobuf_field_mask_pb from 'google-protobuf/google/protobuf/field_mask_pb'; // proto import: "google/protobuf/field_mask.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as types_time_period_pb from '@smart-core-os/sc-api-grpc-web/types/time/period_pb'; // proto import: "types/time/period.proto"


export class EmergencyLight extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EmergencyLight.AsObject;
  static toObject(includeInstance: boolean, msg: EmergencyLight): EmergencyLight.AsObject;
  static serializeBinaryToWriter(message: EmergencyLight, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EmergencyLight;
  static deserializeBinaryFromReader(message: EmergencyLight, reader: jspb.BinaryReader): EmergencyLight;
}

export namespace EmergencyLight {
  export type AsObject = {
  }

  export enum TestType { 
    TEST_UNKNOWN = 0,
    NO_TEST = 1,
    FUNCTION_TEST = 2,
    DURATION_TEST = 3,
  }

  export enum Result { 
    TEST_RESULT_UNSPECIFIED = 0,
    TEST_PASSED = 1,
    CIRCUIT_FAILURE = 2,
    BATTERY_DURATION_FAILURE = 3,
    BATTERY_FAILURE = 4,
    LAMP_FAILURE = 5,
    FUNCTION_TEST_FAILED = 6,
    DURATION_TEST_FAILED = 7,
    LIGHT_FAULTY = 8,
    COMMUNICATION_FAILURE = 9,
    OTHER_FAULT = 10,
  }
}

export class StartTestRequest extends jspb.Message {
  getName(): string;
  setName(value: string): StartTestRequest;

  getTest(): EmergencyLight.TestType;
  setTest(value: EmergencyLight.TestType): StartTestRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartTestRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StartTestRequest): StartTestRequest.AsObject;
  static serializeBinaryToWriter(message: StartTestRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartTestRequest;
  static deserializeBinaryFromReader(message: StartTestRequest, reader: jspb.BinaryReader): StartTestRequest;
}

export namespace StartTestRequest {
  export type AsObject = {
    name: string,
    test: EmergencyLight.TestType,
  }
}

export class StartTestResponse extends jspb.Message {
  getTest(): EmergencyLight.TestType;
  setTest(value: EmergencyLight.TestType): StartTestResponse;

  getStartTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setStartTime(value?: google_protobuf_timestamp_pb.Timestamp): StartTestResponse;
  hasStartTime(): boolean;
  clearStartTime(): StartTestResponse;

  getDuration(): google_protobuf_duration_pb.Duration | undefined;
  setDuration(value?: google_protobuf_duration_pb.Duration): StartTestResponse;
  hasDuration(): boolean;
  clearDuration(): StartTestResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartTestResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StartTestResponse): StartTestResponse.AsObject;
  static serializeBinaryToWriter(message: StartTestResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartTestResponse;
  static deserializeBinaryFromReader(message: StartTestResponse, reader: jspb.BinaryReader): StartTestResponse;
}

export namespace StartTestResponse {
  export type AsObject = {
    test: EmergencyLight.TestType,
    startTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    duration?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class StopTestRequest extends jspb.Message {
  getName(): string;
  setName(value: string): StopTestRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopTestRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StopTestRequest): StopTestRequest.AsObject;
  static serializeBinaryToWriter(message: StopTestRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StopTestRequest;
  static deserializeBinaryFromReader(message: StopTestRequest, reader: jspb.BinaryReader): StopTestRequest;
}

export namespace StopTestRequest {
  export type AsObject = {
    name: string,
  }
}

export class StopTestResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopTestResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StopTestResponse): StopTestResponse.AsObject;
  static serializeBinaryToWriter(message: StopTestResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StopTestResponse;
  static deserializeBinaryFromReader(message: StopTestResponse, reader: jspb.BinaryReader): StopTestResponse;
}

export namespace StopTestResponse {
  export type AsObject = {
  }
}

export class GetTestResultsRequest extends jspb.Message {
  getName(): string;
  setName(value: string): GetTestResultsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetTestResultsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetTestResultsRequest): GetTestResultsRequest.AsObject;
  static serializeBinaryToWriter(message: GetTestResultsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetTestResultsRequest;
  static deserializeBinaryFromReader(message: GetTestResultsRequest, reader: jspb.BinaryReader): GetTestResultsRequest;
}

export namespace GetTestResultsRequest {
  export type AsObject = {
    name: string,
  }
}

export class GetTestResultsResponse extends jspb.Message {
  getFunctionTestResult(): TestResult | undefined;
  setFunctionTestResult(value?: TestResult): GetTestResultsResponse;
  hasFunctionTestResult(): boolean;
  clearFunctionTestResult(): GetTestResultsResponse;

  getDurationTestResult(): TestResult | undefined;
  setDurationTestResult(value?: TestResult): GetTestResultsResponse;
  hasDurationTestResult(): boolean;
  clearDurationTestResult(): GetTestResultsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetTestResultsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetTestResultsResponse): GetTestResultsResponse.AsObject;
  static serializeBinaryToWriter(message: GetTestResultsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetTestResultsResponse;
  static deserializeBinaryFromReader(message: GetTestResultsResponse, reader: jspb.BinaryReader): GetTestResultsResponse;
}

export namespace GetTestResultsResponse {
  export type AsObject = {
    functionTestResult?: TestResult.AsObject,
    durationTestResult?: TestResult.AsObject,
  }
}

export class TestResult extends jspb.Message {
  getTestType(): EmergencyLight.TestType;
  setTestType(value: EmergencyLight.TestType): TestResult;

  getFunctionTestStart(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setFunctionTestStart(value?: google_protobuf_timestamp_pb.Timestamp): TestResult;
  hasFunctionTestStart(): boolean;
  clearFunctionTestStart(): TestResult;

  getFunctionTestCompletion(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setFunctionTestCompletion(value?: google_protobuf_timestamp_pb.Timestamp): TestResult;
  hasFunctionTestCompletion(): boolean;
  clearFunctionTestCompletion(): TestResult;

  getFunctionTestResult(): EmergencyLight.Result;
  setFunctionTestResult(value: EmergencyLight.Result): TestResult;

  getDurationTestStart(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setDurationTestStart(value?: google_protobuf_timestamp_pb.Timestamp): TestResult;
  hasDurationTestStart(): boolean;
  clearDurationTestStart(): TestResult;

  getDurationTestCompletion(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setDurationTestCompletion(value?: google_protobuf_timestamp_pb.Timestamp): TestResult;
  hasDurationTestCompletion(): boolean;
  clearDurationTestCompletion(): TestResult;

  getDurationTestResult(): EmergencyLight.Result;
  setDurationTestResult(value: EmergencyLight.Result): TestResult;

  getDuration(): google_protobuf_duration_pb.Duration | undefined;
  setDuration(value?: google_protobuf_duration_pb.Duration): TestResult;
  hasDuration(): boolean;
  clearDuration(): TestResult;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TestResult.AsObject;
  static toObject(includeInstance: boolean, msg: TestResult): TestResult.AsObject;
  static serializeBinaryToWriter(message: TestResult, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TestResult;
  static deserializeBinaryFromReader(message: TestResult, reader: jspb.BinaryReader): TestResult;
}

export namespace TestResult {
  export type AsObject = {
    testType: EmergencyLight.TestType,
    functionTestStart?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    functionTestCompletion?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    functionTestResult: EmergencyLight.Result,
    durationTestStart?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    durationTestCompletion?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    durationTestResult: EmergencyLight.Result,
    duration?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class ListTestResultsRequest extends jspb.Message {
  getPeriod(): types_time_period_pb.Period | undefined;
  setPeriod(value?: types_time_period_pb.Period): ListTestResultsRequest;
  hasPeriod(): boolean;
  clearPeriod(): ListTestResultsRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): ListTestResultsRequest;
  hasReadMask(): boolean;
  clearReadMask(): ListTestResultsRequest;

  getPageSize(): number;
  setPageSize(value: number): ListTestResultsRequest;

  getPageToken(): string;
  setPageToken(value: string): ListTestResultsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTestResultsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListTestResultsRequest): ListTestResultsRequest.AsObject;
  static serializeBinaryToWriter(message: ListTestResultsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListTestResultsRequest;
  static deserializeBinaryFromReader(message: ListTestResultsRequest, reader: jspb.BinaryReader): ListTestResultsRequest;
}

export namespace ListTestResultsRequest {
  export type AsObject = {
    period?: types_time_period_pb.Period.AsObject,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
    pageSize: number,
    pageToken: string,
  }
}

export class ListTestResultsResponse extends jspb.Message {
  getTestResultsList(): Array<TestResult>;
  setTestResultsList(value: Array<TestResult>): ListTestResultsResponse;
  clearTestResultsList(): ListTestResultsResponse;
  addTestResults(value?: TestResult, index?: number): TestResult;

  getNextPageToken(): string;
  setNextPageToken(value: string): ListTestResultsResponse;

  getTotalSize(): number;
  setTotalSize(value: number): ListTestResultsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTestResultsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListTestResultsResponse): ListTestResultsResponse.AsObject;
  static serializeBinaryToWriter(message: ListTestResultsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListTestResultsResponse;
  static deserializeBinaryFromReader(message: ListTestResultsResponse, reader: jspb.BinaryReader): ListTestResultsResponse;
}

export namespace ListTestResultsResponse {
  export type AsObject = {
    testResultsList: Array<TestResult.AsObject>,
    nextPageToken: string,
    totalSize: number,
  }
}

