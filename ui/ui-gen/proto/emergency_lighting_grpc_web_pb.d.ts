import * as grpcWeb from 'grpc-web';

import * as proto_emergency_lighting_pb from '../proto/emergency_lighting_pb';


export class EmergencyLightingApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getEmergencyLight(
    request: proto_emergency_lighting_pb.GetEmergencyLightRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_emergency_lighting_pb.EmergencyLight) => void
  ): grpcWeb.ClientReadableStream<proto_emergency_lighting_pb.EmergencyLight>;

  listEmergencyLights(
    request: proto_emergency_lighting_pb.ListEmergencyLightsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_emergency_lighting_pb.ListEmergencyLightsResponse) => void
  ): grpcWeb.ClientReadableStream<proto_emergency_lighting_pb.ListEmergencyLightsResponse>;

  listEmergencyLightEvents(
    request: proto_emergency_lighting_pb.ListEmergencyLightEventsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_emergency_lighting_pb.ListEmergencyLightEventsResponse) => void
  ): grpcWeb.ClientReadableStream<proto_emergency_lighting_pb.ListEmergencyLightEventsResponse>;

  getReportCSV(
    request: proto_emergency_lighting_pb.GetReportCSVRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_emergency_lighting_pb.ReportCSV) => void
  ): grpcWeb.ClientReadableStream<proto_emergency_lighting_pb.ReportCSV>;

}

export class EmergencyLightingApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getEmergencyLight(
    request: proto_emergency_lighting_pb.GetEmergencyLightRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_emergency_lighting_pb.EmergencyLight>;

  listEmergencyLights(
    request: proto_emergency_lighting_pb.ListEmergencyLightsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_emergency_lighting_pb.ListEmergencyLightsResponse>;

  listEmergencyLightEvents(
    request: proto_emergency_lighting_pb.ListEmergencyLightEventsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_emergency_lighting_pb.ListEmergencyLightEventsResponse>;

  getReportCSV(
    request: proto_emergency_lighting_pb.GetReportCSVRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_emergency_lighting_pb.ReportCSV>;

}

