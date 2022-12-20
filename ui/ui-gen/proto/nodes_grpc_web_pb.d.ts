import * as grpcWeb from 'grpc-web';

import * as proto_nodes_pb from '../proto/nodes_pb';


export class NodeApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getNodeRegistration(
    request: proto_nodes_pb.GetNodeRegistrationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_nodes_pb.NodeRegistration) => void
  ): grpcWeb.ClientReadableStream<proto_nodes_pb.NodeRegistration>;

  createNodeRegistration(
    request: proto_nodes_pb.CreateNodeRegistrationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_nodes_pb.NodeRegistration) => void
  ): grpcWeb.ClientReadableStream<proto_nodes_pb.NodeRegistration>;

  listNodeRegistrations(
    request: proto_nodes_pb.ListNodeRegistrationsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_nodes_pb.ListNodeRegistrationsResponse) => void
  ): grpcWeb.ClientReadableStream<proto_nodes_pb.ListNodeRegistrationsResponse>;

  testNodeCommunication(
    request: proto_nodes_pb.TestNodeCommunicationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_nodes_pb.TestNodeCommunicationResponse) => void
  ): grpcWeb.ClientReadableStream<proto_nodes_pb.TestNodeCommunicationResponse>;

}

export class NodeApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getNodeRegistration(
    request: proto_nodes_pb.GetNodeRegistrationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_nodes_pb.NodeRegistration>;

  createNodeRegistration(
    request: proto_nodes_pb.CreateNodeRegistrationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_nodes_pb.NodeRegistration>;

  listNodeRegistrations(
    request: proto_nodes_pb.ListNodeRegistrationsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_nodes_pb.ListNodeRegistrationsResponse>;

  testNodeCommunication(
    request: proto_nodes_pb.TestNodeCommunicationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_nodes_pb.TestNodeCommunicationResponse>;

}

