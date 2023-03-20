import * as jspb from 'google-protobuf'



export class HubNode extends jspb.Message {
  getAddress(): string;
  setAddress(value: string): HubNode;

  getName(): string;
  setName(value: string): HubNode;

  getDescription(): string;
  setDescription(value: string): HubNode;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HubNode.AsObject;
  static toObject(includeInstance: boolean, msg: HubNode): HubNode.AsObject;
  static serializeBinaryToWriter(message: HubNode, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HubNode;
  static deserializeBinaryFromReader(message: HubNode, reader: jspb.BinaryReader): HubNode;
}

export namespace HubNode {
  export type AsObject = {
    address: string,
    name: string,
    description: string,
  }
}

export class GetHubNodeRequest extends jspb.Message {
  getAddress(): string;
  setAddress(value: string): GetHubNodeRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetHubNodeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetHubNodeRequest): GetHubNodeRequest.AsObject;
  static serializeBinaryToWriter(message: GetHubNodeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetHubNodeRequest;
  static deserializeBinaryFromReader(message: GetHubNodeRequest, reader: jspb.BinaryReader): GetHubNodeRequest;
}

export namespace GetHubNodeRequest {
  export type AsObject = {
    address: string,
  }
}

export class EnrollHubNodeRequest extends jspb.Message {
  getNode(): HubNode | undefined;
  setNode(value?: HubNode): EnrollHubNodeRequest;
  hasNode(): boolean;
  clearNode(): EnrollHubNodeRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnrollHubNodeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EnrollHubNodeRequest): EnrollHubNodeRequest.AsObject;
  static serializeBinaryToWriter(message: EnrollHubNodeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnrollHubNodeRequest;
  static deserializeBinaryFromReader(message: EnrollHubNodeRequest, reader: jspb.BinaryReader): EnrollHubNodeRequest;
}

export namespace EnrollHubNodeRequest {
  export type AsObject = {
    node?: HubNode.AsObject,
  }
}

export class ListHubNodesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListHubNodesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListHubNodesRequest): ListHubNodesRequest.AsObject;
  static serializeBinaryToWriter(message: ListHubNodesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListHubNodesRequest;
  static deserializeBinaryFromReader(message: ListHubNodesRequest, reader: jspb.BinaryReader): ListHubNodesRequest;
}

export namespace ListHubNodesRequest {
  export type AsObject = {
  }
}

export class ListHubNodesResponse extends jspb.Message {
  getNodesList(): Array<HubNode>;
  setNodesList(value: Array<HubNode>): ListHubNodesResponse;
  clearNodesList(): ListHubNodesResponse;
  addNodes(value?: HubNode, index?: number): HubNode;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListHubNodesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListHubNodesResponse): ListHubNodesResponse.AsObject;
  static serializeBinaryToWriter(message: ListHubNodesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListHubNodesResponse;
  static deserializeBinaryFromReader(message: ListHubNodesResponse, reader: jspb.BinaryReader): ListHubNodesResponse;
}

export namespace ListHubNodesResponse {
  export type AsObject = {
    nodesList: Array<HubNode.AsObject>,
  }
}

export class TestHubNodeRequest extends jspb.Message {
  getAddress(): string;
  setAddress(value: string): TestHubNodeRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TestHubNodeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: TestHubNodeRequest): TestHubNodeRequest.AsObject;
  static serializeBinaryToWriter(message: TestHubNodeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TestHubNodeRequest;
  static deserializeBinaryFromReader(message: TestHubNodeRequest, reader: jspb.BinaryReader): TestHubNodeRequest;
}

export namespace TestHubNodeRequest {
  export type AsObject = {
    address: string,
  }
}

export class TestHubNodeResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TestHubNodeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: TestHubNodeResponse): TestHubNodeResponse.AsObject;
  static serializeBinaryToWriter(message: TestHubNodeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TestHubNodeResponse;
  static deserializeBinaryFromReader(message: TestHubNodeResponse, reader: jspb.BinaryReader): TestHubNodeResponse;
}

export namespace TestHubNodeResponse {
  export type AsObject = {
  }
}

