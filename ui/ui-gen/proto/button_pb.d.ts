import * as jspb from 'google-protobuf'

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';


export class ButtonState extends jspb.Message {
  getState(): ButtonState.Press;
  setState(value: ButtonState.Press): ButtonState;

  getStateChangeTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setStateChangeTime(value?: google_protobuf_timestamp_pb.Timestamp): ButtonState;
  hasStateChangeTime(): boolean;
  clearStateChangeTime(): ButtonState;

  getMostRecentGesture(): ButtonState.Gesture | undefined;
  setMostRecentGesture(value?: ButtonState.Gesture): ButtonState;
  hasMostRecentGesture(): boolean;
  clearMostRecentGesture(): ButtonState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ButtonState.AsObject;
  static toObject(includeInstance: boolean, msg: ButtonState): ButtonState.AsObject;
  static serializeBinaryToWriter(message: ButtonState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ButtonState;
  static deserializeBinaryFromReader(message: ButtonState, reader: jspb.BinaryReader): ButtonState;
}

export namespace ButtonState {
  export type AsObject = {
    state: ButtonState.Press,
    stateChangeTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    mostRecentGesture?: ButtonState.Gesture.AsObject,
  }

  export class Gesture extends jspb.Message {
    getId(): string;
    setId(value: string): Gesture;

    getKind(): ButtonState.Gesture.Kind;
    setKind(value: ButtonState.Gesture.Kind): Gesture;

    getCount(): number;
    setCount(value: number): Gesture;

    getStartTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setStartTime(value?: google_protobuf_timestamp_pb.Timestamp): Gesture;
    hasStartTime(): boolean;
    clearStartTime(): Gesture;

    getEndTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setEndTime(value?: google_protobuf_timestamp_pb.Timestamp): Gesture;
    hasEndTime(): boolean;
    clearEndTime(): Gesture;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Gesture.AsObject;
    static toObject(includeInstance: boolean, msg: Gesture): Gesture.AsObject;
    static serializeBinaryToWriter(message: Gesture, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Gesture;
    static deserializeBinaryFromReader(message: Gesture, reader: jspb.BinaryReader): Gesture;
  }

  export namespace Gesture {
    export type AsObject = {
      id: string,
      kind: ButtonState.Gesture.Kind,
      count: number,
      startTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
      endTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    }

    export enum Kind { 
      KIND_UNSPECIFIED = 0,
      CLICK = 1,
      HOLD = 2,
    }
  }


  export enum Press { 
    STATE_UNSPECIFIED = 0,
    UNPRESSED = 1,
    PRESSED = 2,
  }
}

export class GetButtonStateRequest extends jspb.Message {
  getName(): string;
  setName(value: string): GetButtonStateRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetButtonStateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetButtonStateRequest): GetButtonStateRequest.AsObject;
  static serializeBinaryToWriter(message: GetButtonStateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetButtonStateRequest;
  static deserializeBinaryFromReader(message: GetButtonStateRequest, reader: jspb.BinaryReader): GetButtonStateRequest;
}

export namespace GetButtonStateRequest {
  export type AsObject = {
    name: string,
  }
}

export class PullButtonStateRequest extends jspb.Message {
  getName(): string;
  setName(value: string): PullButtonStateRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullButtonStateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PullButtonStateRequest): PullButtonStateRequest.AsObject;
  static serializeBinaryToWriter(message: PullButtonStateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullButtonStateRequest;
  static deserializeBinaryFromReader(message: PullButtonStateRequest, reader: jspb.BinaryReader): PullButtonStateRequest;
}

export namespace PullButtonStateRequest {
  export type AsObject = {
    name: string,
  }
}

export class PullButtonStateResponse extends jspb.Message {
  getChangesList(): Array<PullButtonStateResponse.Change>;
  setChangesList(value: Array<PullButtonStateResponse.Change>): PullButtonStateResponse;
  clearChangesList(): PullButtonStateResponse;
  addChanges(value?: PullButtonStateResponse.Change, index?: number): PullButtonStateResponse.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullButtonStateResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PullButtonStateResponse): PullButtonStateResponse.AsObject;
  static serializeBinaryToWriter(message: PullButtonStateResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullButtonStateResponse;
  static deserializeBinaryFromReader(message: PullButtonStateResponse, reader: jspb.BinaryReader): PullButtonStateResponse;
}

export namespace PullButtonStateResponse {
  export type AsObject = {
    changesList: Array<PullButtonStateResponse.Change.AsObject>,
  }

  export class Change extends jspb.Message {
    getName(): string;
    setName(value: string): Change;

    getChangeTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setChangeTime(value?: google_protobuf_timestamp_pb.Timestamp): Change;
    hasChangeTime(): boolean;
    clearChangeTime(): Change;

    getButtonState(): ButtonState | undefined;
    setButtonState(value?: ButtonState): Change;
    hasButtonState(): boolean;
    clearButtonState(): Change;

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
      buttonState?: ButtonState.AsObject,
    }
  }

}

