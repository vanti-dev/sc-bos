import * as jspb from 'google-protobuf'



export class ClearPriorityValueRequest extends jspb.Message {
  getName(): string;
  setName(value: string): ClearPriorityValueRequest;

  getTrait(): string;
  setTrait(value: string): ClearPriorityValueRequest;

  getEntryIndex(): number;
  setEntryIndex(value: number): ClearPriorityValueRequest;

  getEntryName(): string;
  setEntryName(value: string): ClearPriorityValueRequest;

  getIdCase(): ClearPriorityValueRequest.IdCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClearPriorityValueRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ClearPriorityValueRequest): ClearPriorityValueRequest.AsObject;
  static serializeBinaryToWriter(message: ClearPriorityValueRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClearPriorityValueRequest;
  static deserializeBinaryFromReader(message: ClearPriorityValueRequest, reader: jspb.BinaryReader): ClearPriorityValueRequest;
}

export namespace ClearPriorityValueRequest {
  export type AsObject = {
    name: string,
    trait: string,
    entryIndex: number,
    entryName: string,
  }

  export enum IdCase { 
    ID_NOT_SET = 0,
    ENTRY_INDEX = 3,
    ENTRY_NAME = 4,
  }
}

export class ClearPriorityValueResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClearPriorityValueResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ClearPriorityValueResponse): ClearPriorityValueResponse.AsObject;
  static serializeBinaryToWriter(message: ClearPriorityValueResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClearPriorityValueResponse;
  static deserializeBinaryFromReader(message: ClearPriorityValueResponse, reader: jspb.BinaryReader): ClearPriorityValueResponse;
}

export namespace ClearPriorityValueResponse {
  export type AsObject = {
  }
}

