import * as jspb from 'google-protobuf'



export class Priority extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): Priority;

  getLevel(): Priority.Level;
  setLevel(value: Priority.Level): Priority;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Priority.AsObject;
  static toObject(includeInstance: boolean, msg: Priority): Priority.AsObject;
  static serializeBinaryToWriter(message: Priority, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Priority;
  static deserializeBinaryFromReader(message: Priority, reader: jspb.BinaryReader): Priority;
}

export namespace Priority {
  export type AsObject = {
    title: string,
    level: Priority.Level,
  }

  export enum Level { 
    LEVEL_UNSPECIFIED = 0,
    MANUAL_LIFE_SAFETY = 10,
    AUTOMATIC_LIFE_SAFETY = 20,
    CRITICAL_OVERRIDE = 30,
    DEBOUNCE = 40,
    SUPERVISOR = 50,
    OPERATOR = 60,
    USER = 70,
    AUTOMATION = 100,
    DEFAULT = 150,
    MINIMUM = 255,
  }
}

