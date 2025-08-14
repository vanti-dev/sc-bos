import * as grpcWeb from 'grpc-web';

import * as emergency_light_pb from './emergency_light_pb'; // proto import: "emergency_light.proto"


export class EmergencyLightApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  startFunctionTest(
    request: emergency_light_pb.StartEmergencyTestRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: emergency_light_pb.StartEmergencyTestResponse) => void
  ): grpcWeb.ClientReadableStream<emergency_light_pb.StartEmergencyTestResponse>;

  startDurationTest(
    request: emergency_light_pb.StartEmergencyTestRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: emergency_light_pb.StartEmergencyTestResponse) => void
  ): grpcWeb.ClientReadableStream<emergency_light_pb.StartEmergencyTestResponse>;

  stopEmergencyTest(
    request: emergency_light_pb.StopEmergencyTestsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: emergency_light_pb.StopEmergencyTestsResponse) => void
  ): grpcWeb.ClientReadableStream<emergency_light_pb.StopEmergencyTestsResponse>;

  getTestResultSet(
    request: emergency_light_pb.GetTestResultSetRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: emergency_light_pb.TestResultSet) => void
  ): grpcWeb.ClientReadableStream<emergency_light_pb.TestResultSet>;

  pullTestResultSets(
    request: emergency_light_pb.PullTestResultRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<emergency_light_pb.PullTestResultsResponse>;

}

export class EmergencyLightApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  startFunctionTest(
    request: emergency_light_pb.StartEmergencyTestRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<emergency_light_pb.StartEmergencyTestResponse>;

  startDurationTest(
    request: emergency_light_pb.StartEmergencyTestRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<emergency_light_pb.StartEmergencyTestResponse>;

  stopEmergencyTest(
    request: emergency_light_pb.StopEmergencyTestsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<emergency_light_pb.StopEmergencyTestsResponse>;

  getTestResultSet(
    request: emergency_light_pb.GetTestResultSetRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<emergency_light_pb.TestResultSet>;

  pullTestResultSets(
    request: emergency_light_pb.PullTestResultRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<emergency_light_pb.PullTestResultsResponse>;

}

