import * as grpcWeb from 'grpc-web';

import * as hub_pb from './hub_pb'; // proto import: "hub.proto"


export class HubApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getHubNode(
    request: hub_pb.GetHubNodeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: hub_pb.HubNode) => void
  ): grpcWeb.ClientReadableStream<hub_pb.HubNode>;

  listHubNodes(
    request: hub_pb.ListHubNodesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: hub_pb.ListHubNodesResponse) => void
  ): grpcWeb.ClientReadableStream<hub_pb.ListHubNodesResponse>;

  pullHubNodes(
    request: hub_pb.PullHubNodesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<hub_pb.PullHubNodesResponse>;

  inspectHubNode(
    request: hub_pb.InspectHubNodeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: hub_pb.HubNodeInspection) => void
  ): grpcWeb.ClientReadableStream<hub_pb.HubNodeInspection>;

  enrollHubNode(
    request: hub_pb.EnrollHubNodeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: hub_pb.HubNode) => void
  ): grpcWeb.ClientReadableStream<hub_pb.HubNode>;

  renewHubNode(
    request: hub_pb.RenewHubNodeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: hub_pb.HubNode) => void
  ): grpcWeb.ClientReadableStream<hub_pb.HubNode>;

  testHubNode(
    request: hub_pb.TestHubNodeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: hub_pb.TestHubNodeResponse) => void
  ): grpcWeb.ClientReadableStream<hub_pb.TestHubNodeResponse>;

  forgetHubNode(
    request: hub_pb.ForgetHubNodeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: hub_pb.ForgetHubNodeResponse) => void
  ): grpcWeb.ClientReadableStream<hub_pb.ForgetHubNodeResponse>;

}

export class HubApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getHubNode(
    request: hub_pb.GetHubNodeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<hub_pb.HubNode>;

  listHubNodes(
    request: hub_pb.ListHubNodesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<hub_pb.ListHubNodesResponse>;

  pullHubNodes(
    request: hub_pb.PullHubNodesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<hub_pb.PullHubNodesResponse>;

  inspectHubNode(
    request: hub_pb.InspectHubNodeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<hub_pb.HubNodeInspection>;

  enrollHubNode(
    request: hub_pb.EnrollHubNodeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<hub_pb.HubNode>;

  renewHubNode(
    request: hub_pb.RenewHubNodeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<hub_pb.HubNode>;

  testHubNode(
    request: hub_pb.TestHubNodeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<hub_pb.TestHubNodeResponse>;

  forgetHubNode(
    request: hub_pb.ForgetHubNodeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<hub_pb.ForgetHubNodeResponse>;

}

