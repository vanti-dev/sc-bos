import * as jspb from 'google-protobuf'



export class Location extends jspb.Message {
  getId(): string;
  setId(value: string): Location;

  getTitle(): string;
  setTitle(value: string): Location;

  getDescription(): string;
  setDescription(value: string): Location;

  getFloor(): string;
  setFloor(value: string): Location;

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
    title: string,
    description: string,
    floor: string,
  }
}

