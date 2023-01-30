import * as jspb from 'google-protobuf'

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';


export class ButtonState extends jspb.Message {
  getPressState(): ButtonState.PressState;
  setPressState(value: ButtonState.PressState): ButtonState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ButtonState.AsObject;
  static toObject(includeInstance: boolean, msg: ButtonState): ButtonState.AsObject;
  static serializeBinaryToWriter(message: ButtonState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ButtonState;
  static deserializeBinaryFromReader(message: ButtonState, reader: jspb.BinaryReader): ButtonState;
}

export namespace ButtonState {
  export type AsObject = {
    pressState: ButtonState.PressState,
  }

  export enum PressState { 
    BUTTON_STATE_UNSPECIFIED = 0,
    PRESSED = 1,
    RELEASED = 2,
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

export class PullButtonEventsRequest extends jspb.Message {
  getName(): string;
  setName(value: string): PullButtonEventsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullButtonEventsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PullButtonEventsRequest): PullButtonEventsRequest.AsObject;
  static serializeBinaryToWriter(message: PullButtonEventsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullButtonEventsRequest;
  static deserializeBinaryFromReader(message: PullButtonEventsRequest, reader: jspb.BinaryReader): PullButtonEventsRequest;
}

export namespace PullButtonEventsRequest {
  export type AsObject = {
    name: string,
  }
}

export class PullButtonEventsResponse extends jspb.Message {
  getChangesList(): Array<PullButtonEventsResponse.Change>;
  setChangesList(value: Array<PullButtonEventsResponse.Change>): PullButtonEventsResponse;
  clearChangesList(): PullButtonEventsResponse;
  addChanges(value?: PullButtonEventsResponse.Change, index?: number): PullButtonEventsResponse.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullButtonEventsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PullButtonEventsResponse): PullButtonEventsResponse.AsObject;
  static serializeBinaryToWriter(message: PullButtonEventsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullButtonEventsResponse;
  static deserializeBinaryFromReader(message: PullButtonEventsResponse, reader: jspb.BinaryReader): PullButtonEventsResponse;
}

export namespace PullButtonEventsResponse {
  export type AsObject = {
    changesList: Array<PullButtonEventsResponse.Change.AsObject>,
  }

  export class Change extends jspb.Message {
    getName(): string;
    setName(value: string): Change;

    getChangeTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setChangeTime(value?: google_protobuf_timestamp_pb.Timestamp): Change;
    hasChangeTime(): boolean;
    clearChangeTime(): Change;

    getButtonEvent(): ButtonEvent;
    setButtonEvent(value: ButtonEvent): Change;

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
      buttonEvent: ButtonEvent,
    }
  }

}

export enum ButtonEvent { 
  BUTTON_EVENT_UNSPECIFIED = 0,
  PRESS = 1,
  RELEASE = 2,
  SHORT_PRESS = 3,
  DOUBLE_PRESS = 4,
  LONG_PRESS_START = 5,
  LONG_PRESS_REPEAT = 6,
  LONG_PRESS_END = 7,
}
