import * as grpcWeb from 'grpc-web';

import * as nodes_pb from './nodes_pb';


export class NodeApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getNodeRegistration(
    request: nodes_pb.GetNodeRegistrationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: nodes_pb.NodeRegistration) => void
  ): grpcWeb.ClientReadableStream<nodes_pb.NodeRegistration>;

  createNodeRegistration(
    request: nodes_pb.CreateNodeRegistrationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: nodes_pb.NodeRegistration) => void
  ): grpcWeb.ClientReadableStream<nodes_pb.NodeRegistration>;

  listNodeRegistrations(
    request: nodes_pb.ListNodeRegistrationsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: nodes_pb.ListNodeRegistrationsResponse) => void
  ): grpcWeb.ClientReadableStream<nodes_pb.ListNodeRegistrationsResponse>;

  testNodeCommunication(
    request: nodes_pb.TestNodeCommunicationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: nodes_pb.TestNodeCommunicationResponse) => void
  ): grpcWeb.ClientReadableStream<nodes_pb.TestNodeCommunicationResponse>;

}

export class NodeApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getNodeRegistration(
    request: nodes_pb.GetNodeRegistrationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<nodes_pb.NodeRegistration>;

  createNodeRegistration(
    request: nodes_pb.CreateNodeRegistrationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<nodes_pb.NodeRegistration>;

  listNodeRegistrations(
    request: nodes_pb.ListNodeRegistrationsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<nodes_pb.ListNodeRegistrationsResponse>;

  testNodeCommunication(
    request: nodes_pb.TestNodeCommunicationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<nodes_pb.TestNodeCommunicationResponse>;

}

