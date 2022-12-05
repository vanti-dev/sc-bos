import * as jspb from 'google-protobuf'


export class Enrollment extends jspb.Message {
  getTargetName(): string;
  setTargetName(value: string): Enrollment;

  getTargetAddress(): string;
  setTargetAddress(value: string): Enrollment;

  getManagerName(): string;
  setManagerName(value: string): Enrollment;

  getManagerAddress(): string;
  setManagerAddress(value: string): Enrollment;

  getCertificate(): Uint8Array | string;
  getCertificate_asU8(): Uint8Array;
  getCertificate_asB64(): string;
  setCertificate(value: Uint8Array | string): Enrollment;

  getRootCas(): Uint8Array | string;
  getRootCas_asU8(): Uint8Array;
  getRootCas_asB64(): string;
  setRootCas(value: Uint8Array | string): Enrollment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Enrollment.AsObject;
  static toObject(includeInstance: boolean, msg: Enrollment): Enrollment.AsObject;
  static serializeBinaryToWriter(message: Enrollment, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Enrollment;
  static deserializeBinaryFromReader(message: Enrollment, reader: jspb.BinaryReader): Enrollment;
}

export namespace Enrollment {
  export type AsObject = {
    targetName: string,
    targetAddress: string,
    managerName: string,
    managerAddress: string,
    certificate: Uint8Array | string,
    rootCas: Uint8Array | string,
  }
}

export class GetEnrollmentRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetEnrollmentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetEnrollmentRequest): GetEnrollmentRequest.AsObject;
  static serializeBinaryToWriter(message: GetEnrollmentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetEnrollmentRequest;
  static deserializeBinaryFromReader(message: GetEnrollmentRequest, reader: jspb.BinaryReader): GetEnrollmentRequest;
}

export namespace GetEnrollmentRequest {
  export type AsObject = {}
}

export class CreateEnrollmentRequest extends jspb.Message {
  getEnrollment(): Enrollment | undefined;
  setEnrollment(value?: Enrollment): CreateEnrollmentRequest;
  hasEnrollment(): boolean;
  clearEnrollment(): CreateEnrollmentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateEnrollmentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateEnrollmentRequest): CreateEnrollmentRequest.AsObject;
  static serializeBinaryToWriter(message: CreateEnrollmentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateEnrollmentRequest;
  static deserializeBinaryFromReader(message: CreateEnrollmentRequest, reader: jspb.BinaryReader): CreateEnrollmentRequest;
}

export namespace CreateEnrollmentRequest {
  export type AsObject = {
    enrollment?: Enrollment.AsObject,
  }
}

export class DeleteEnrollmentRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteEnrollmentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteEnrollmentRequest): DeleteEnrollmentRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteEnrollmentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteEnrollmentRequest;
  static deserializeBinaryFromReader(message: DeleteEnrollmentRequest, reader: jspb.BinaryReader): DeleteEnrollmentRequest;
}

export namespace DeleteEnrollmentRequest {
  export type AsObject = {}
}

