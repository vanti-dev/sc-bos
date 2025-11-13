import * as grpcWeb from 'grpc-web';

import * as fluid_flow_pb from './fluid_flow_pb'; // proto import: "fluid_flow.proto"


export class FluidFlowApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getFluidFlow(
    request: fluid_flow_pb.GetFluidFlowRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: fluid_flow_pb.FluidFlow) => void
  ): grpcWeb.ClientReadableStream<fluid_flow_pb.FluidFlow>;

  pullFluidFlow(
    request: fluid_flow_pb.PullFluidFlowRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<fluid_flow_pb.PullFluidFlowResponse>;

  updateFluidFlow(
    request: fluid_flow_pb.UpdateFluidFlowRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: fluid_flow_pb.FluidFlow) => void
  ): grpcWeb.ClientReadableStream<fluid_flow_pb.FluidFlow>;

}

export class FluidFlowInfoClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describeFluidFlow(
    request: fluid_flow_pb.DescribeFluidFlowRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: fluid_flow_pb.FluidFlowSupport) => void
  ): grpcWeb.ClientReadableStream<fluid_flow_pb.FluidFlowSupport>;

}

export class FluidFlowApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getFluidFlow(
    request: fluid_flow_pb.GetFluidFlowRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<fluid_flow_pb.FluidFlow>;

  pullFluidFlow(
    request: fluid_flow_pb.PullFluidFlowRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<fluid_flow_pb.PullFluidFlowResponse>;

  updateFluidFlow(
    request: fluid_flow_pb.UpdateFluidFlowRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<fluid_flow_pb.FluidFlow>;

}

export class FluidFlowInfoPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describeFluidFlow(
    request: fluid_flow_pb.DescribeFluidFlowRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<fluid_flow_pb.FluidFlowSupport>;

}

