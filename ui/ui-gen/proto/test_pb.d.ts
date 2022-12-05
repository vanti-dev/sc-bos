import * as jspb from 'google-protobuf'


export class Test extends jspb.Message {
  getData(): string;
  setData(value: string): Test;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Test.AsObject;
  static toObject(includeInstance: boolean, msg: Test): Test.AsObject;
  static serializeBinaryToWriter(message: Test, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Test;
  static deserializeBinaryFromReader(message: Test, reader: jspb.BinaryReader): Test;
}

export namespace Test {
  export type AsObject = {
    data: string,
  }
}

export class GetTestRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetTestRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetTestRequest): GetTestRequest.AsObject;
  static serializeBinaryToWriter(message: GetTestRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetTestRequest;
  static deserializeBinaryFromReader(message: GetTestRequest, reader: jspb.BinaryReader): GetTestRequest;
}

export namespace GetTestRequest {
  export type AsObject = {
  }
}

export class UpdateTestRequest extends jspb.Message {
  getTest(): Test | undefined;
  setTest(value?: Test): UpdateTestRequest;
  hasTest(): boolean;
  clearTest(): UpdateTestRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateTestRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateTestRequest): UpdateTestRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateTestRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateTestRequest;
  static deserializeBinaryFromReader(message: UpdateTestRequest, reader: jspb.BinaryReader): UpdateTestRequest;
}

export namespace UpdateTestRequest {
  export type AsObject = {
    test?: Test.AsObject,
  }
}

