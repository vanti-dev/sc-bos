import * as jspb from 'google-protobuf'


export class NodeRegistration extends jspb.Message {
  getName(): string;

  setName(value: string): NodeRegistration;

  getAddress(): string;

  setAddress(value: string): NodeRegistration;

  getDescription(): string;

  setDescription(value: string): NodeRegistration;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): NodeRegistration.AsObject;

  static toObject(includeInstance: boolean, msg: NodeRegistration): NodeRegistration.AsObject;

  static serializeBinaryToWriter(message: NodeRegistration, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): NodeRegistration;

  static deserializeBinaryFromReader(message: NodeRegistration, reader: jspb.BinaryReader): NodeRegistration;
}

export namespace NodeRegistration {
  export type AsObject = {
    name: string,
    address: string,
    description: string,
  }
}

export class GetNodeRegistrationRequest extends jspb.Message {
  getNodeName(): string;

  setNodeName(value: string): GetNodeRegistrationRequest;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): GetNodeRegistrationRequest.AsObject;

  static toObject(includeInstance: boolean, msg: GetNodeRegistrationRequest): GetNodeRegistrationRequest.AsObject;

  static serializeBinaryToWriter(message: GetNodeRegistrationRequest, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): GetNodeRegistrationRequest;

  static deserializeBinaryFromReader(message: GetNodeRegistrationRequest, reader: jspb.BinaryReader): GetNodeRegistrationRequest;
}

export namespace GetNodeRegistrationRequest {
  export type AsObject = {
    nodeName: string,
  }
}

export class CreateNodeRegistrationRequest extends jspb.Message {
  getNodeRegistration(): NodeRegistration | undefined;

  setNodeRegistration(value?: NodeRegistration): CreateNodeRegistrationRequest;

  hasNodeRegistration(): boolean;

  clearNodeRegistration(): CreateNodeRegistrationRequest;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): CreateNodeRegistrationRequest.AsObject;

  static toObject(includeInstance: boolean, msg: CreateNodeRegistrationRequest): CreateNodeRegistrationRequest.AsObject;

  static serializeBinaryToWriter(message: CreateNodeRegistrationRequest, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): CreateNodeRegistrationRequest;

  static deserializeBinaryFromReader(message: CreateNodeRegistrationRequest, reader: jspb.BinaryReader): CreateNodeRegistrationRequest;
}

export namespace CreateNodeRegistrationRequest {
  export type AsObject = {
    nodeRegistration?: NodeRegistration.AsObject,
  }
}

export class ListNodeRegistrationsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): ListNodeRegistrationsRequest.AsObject;

  static toObject(includeInstance: boolean, msg: ListNodeRegistrationsRequest): ListNodeRegistrationsRequest.AsObject;

  static serializeBinaryToWriter(message: ListNodeRegistrationsRequest, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): ListNodeRegistrationsRequest;

  static deserializeBinaryFromReader(message: ListNodeRegistrationsRequest, reader: jspb.BinaryReader): ListNodeRegistrationsRequest;
}

export namespace ListNodeRegistrationsRequest {
  export type AsObject = {}
}

export class ListNodeRegistrationsResponse extends jspb.Message {
  getNodeRegistrationsList(): Array<NodeRegistration>;

  setNodeRegistrationsList(value: Array<NodeRegistration>): ListNodeRegistrationsResponse;

  clearNodeRegistrationsList(): ListNodeRegistrationsResponse;

  addNodeRegistrations(value?: NodeRegistration, index?: number): NodeRegistration;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): ListNodeRegistrationsResponse.AsObject;

  static toObject(includeInstance: boolean, msg: ListNodeRegistrationsResponse): ListNodeRegistrationsResponse.AsObject;

  static serializeBinaryToWriter(message: ListNodeRegistrationsResponse, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): ListNodeRegistrationsResponse;

  static deserializeBinaryFromReader(message: ListNodeRegistrationsResponse, reader: jspb.BinaryReader): ListNodeRegistrationsResponse;
}

export namespace ListNodeRegistrationsResponse {
  export type AsObject = {
    nodeRegistrationsList: Array<NodeRegistration.AsObject>,
  }
}

export class TestNodeCommunicationRequest extends jspb.Message {
  getNodeName(): string;

  setNodeName(value: string): TestNodeCommunicationRequest;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): TestNodeCommunicationRequest.AsObject;

  static toObject(includeInstance: boolean, msg: TestNodeCommunicationRequest): TestNodeCommunicationRequest.AsObject;

  static serializeBinaryToWriter(message: TestNodeCommunicationRequest, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): TestNodeCommunicationRequest;

  static deserializeBinaryFromReader(message: TestNodeCommunicationRequest, reader: jspb.BinaryReader): TestNodeCommunicationRequest;
}

export namespace TestNodeCommunicationRequest {
  export type AsObject = {
    nodeName: string,
  }
}

export class TestNodeCommunicationResponse extends jspb.Message {
  getServicesList(): Array<string>;

  setServicesList(value: Array<string>): TestNodeCommunicationResponse;

  clearServicesList(): TestNodeCommunicationResponse;

  addServices(value: string, index?: number): TestNodeCommunicationResponse;

  serializeBinary(): Uint8Array;

  toObject(includeInstance?: boolean): TestNodeCommunicationResponse.AsObject;

  static toObject(includeInstance: boolean, msg: TestNodeCommunicationResponse): TestNodeCommunicationResponse.AsObject;

  static serializeBinaryToWriter(message: TestNodeCommunicationResponse, writer: jspb.BinaryWriter): void;

  static deserializeBinary(bytes: Uint8Array): TestNodeCommunicationResponse;

  static deserializeBinaryFromReader(message: TestNodeCommunicationResponse, reader: jspb.BinaryReader): TestNodeCommunicationResponse;
}

export namespace TestNodeCommunicationResponse {
  export type AsObject = {
    servicesList: Array<string>,
  }
}

