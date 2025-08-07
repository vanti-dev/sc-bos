import * as grpcWeb from 'grpc-web';

import * as emergency_light_pb from './emergency_light_pb'; // proto import: "emergency_light.proto"


export class EmergencyLightApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  startTest(
    request: emergency_light_pb.StartTestRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: emergency_light_pb.StartTestResponse) => void
  ): grpcWeb.ClientReadableStream<emergency_light_pb.StartTestResponse>;

  stopTest(
    request: emergency_light_pb.StopTestRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: emergency_light_pb.StopTestResponse) => void
  ): grpcWeb.ClientReadableStream<emergency_light_pb.StopTestResponse>;

  getTestResults(
    request: emergency_light_pb.GetTestResultsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: emergency_light_pb.GetTestResultsResponse) => void
  ): grpcWeb.ClientReadableStream<emergency_light_pb.GetTestResultsResponse>;

  listTestResults(
    request: emergency_light_pb.ListTestResultsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: emergency_light_pb.ListTestResultsResponse) => void
  ): grpcWeb.ClientReadableStream<emergency_light_pb.ListTestResultsResponse>;

}

export class EmergencyLightApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  startTest(
    request: emergency_light_pb.StartTestRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<emergency_light_pb.StartTestResponse>;

  stopTest(
    request: emergency_light_pb.StopTestRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<emergency_light_pb.StopTestResponse>;

  getTestResults(
    request: emergency_light_pb.GetTestResultsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<emergency_light_pb.GetTestResultsResponse>;

  listTestResults(
    request: emergency_light_pb.ListTestResultsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<emergency_light_pb.ListTestResultsResponse>;

}

