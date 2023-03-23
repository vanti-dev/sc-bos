import * as grpcWeb from 'grpc-web';

import * as dali_pb from './dali_pb';


export class DaliApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  addToGroup(
    request: dali_pb.AddToGroupRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: dali_pb.AddToGroupResponse) => void
  ): grpcWeb.ClientReadableStream<dali_pb.AddToGroupResponse>;

  removeFromGroup(
    request: dali_pb.RemoveFromGroupRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: dali_pb.RemoveFromGroupResponse) => void
  ): grpcWeb.ClientReadableStream<dali_pb.RemoveFromGroupResponse>;

  getGroupMembership(
    request: dali_pb.GetGroupMembershipRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: dali_pb.GetGroupMembershipResponse) => void
  ): grpcWeb.ClientReadableStream<dali_pb.GetGroupMembershipResponse>;

  getControlGearStatus(
    request: dali_pb.GetControlGearStatusRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: dali_pb.ControlGearStatus) => void
  ): grpcWeb.ClientReadableStream<dali_pb.ControlGearStatus>;

  getEmergencyStatus(
    request: dali_pb.GetEmergencyStatusRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: dali_pb.EmergencyStatus) => void
  ): grpcWeb.ClientReadableStream<dali_pb.EmergencyStatus>;

  identify(
    request: dali_pb.IdentifyRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: dali_pb.IdentifyResponse) => void
  ): grpcWeb.ClientReadableStream<dali_pb.IdentifyResponse>;

  startTest(
    request: dali_pb.StartTestRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: dali_pb.StartTestResponse) => void
  ): grpcWeb.ClientReadableStream<dali_pb.StartTestResponse>;

  stopTest(
    request: dali_pb.StopTestRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: dali_pb.StopTestResponse) => void
  ): grpcWeb.ClientReadableStream<dali_pb.StopTestResponse>;

  getTestResult(
    request: dali_pb.GetTestResultRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: dali_pb.TestResult) => void
  ): grpcWeb.ClientReadableStream<dali_pb.TestResult>;

  deleteTestResult(
    request: dali_pb.DeleteTestResultRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: dali_pb.TestResult) => void
  ): grpcWeb.ClientReadableStream<dali_pb.TestResult>;

}

export class DaliApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  addToGroup(
    request: dali_pb.AddToGroupRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<dali_pb.AddToGroupResponse>;

  removeFromGroup(
    request: dali_pb.RemoveFromGroupRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<dali_pb.RemoveFromGroupResponse>;

  getGroupMembership(
    request: dali_pb.GetGroupMembershipRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<dali_pb.GetGroupMembershipResponse>;

  getControlGearStatus(
    request: dali_pb.GetControlGearStatusRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<dali_pb.ControlGearStatus>;

  getEmergencyStatus(
    request: dali_pb.GetEmergencyStatusRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<dali_pb.EmergencyStatus>;

  identify(
    request: dali_pb.IdentifyRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<dali_pb.IdentifyResponse>;

  startTest(
    request: dali_pb.StartTestRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<dali_pb.StartTestResponse>;

  stopTest(
    request: dali_pb.StopTestRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<dali_pb.StopTestResponse>;

  getTestResult(
    request: dali_pb.GetTestResultRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<dali_pb.TestResult>;

  deleteTestResult(
    request: dali_pb.DeleteTestResultRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<dali_pb.TestResult>;

}

