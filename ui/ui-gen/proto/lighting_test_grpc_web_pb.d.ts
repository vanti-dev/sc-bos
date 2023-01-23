import * as grpcWeb from 'grpc-web';

import * as proto_lighting_test_pb from '../proto/lighting_test_pb';


export class LightingTestApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getLightHealth(
    request: proto_lighting_test_pb.GetLightHealthRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_lighting_test_pb.LightHealth) => void
  ): grpcWeb.ClientReadableStream<proto_lighting_test_pb.LightHealth>;

  listLightHealth(
    request: proto_lighting_test_pb.ListLightHealthRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_lighting_test_pb.ListLightHealthResponse) => void
  ): grpcWeb.ClientReadableStream<proto_lighting_test_pb.ListLightHealthResponse>;

  listLightEvents(
    request: proto_lighting_test_pb.ListLightEventsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_lighting_test_pb.ListLightEventsResponse) => void
  ): grpcWeb.ClientReadableStream<proto_lighting_test_pb.ListLightEventsResponse>;

  getReportCSV(
    request: proto_lighting_test_pb.GetReportCSVRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_lighting_test_pb.ReportCSV) => void
  ): grpcWeb.ClientReadableStream<proto_lighting_test_pb.ReportCSV>;

}

export class LightingTestApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getLightHealth(
    request: proto_lighting_test_pb.GetLightHealthRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_lighting_test_pb.LightHealth>;

  listLightHealth(
    request: proto_lighting_test_pb.ListLightHealthRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_lighting_test_pb.ListLightHealthResponse>;

  listLightEvents(
    request: proto_lighting_test_pb.ListLightEventsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_lighting_test_pb.ListLightEventsResponse>;

  getReportCSV(
    request: proto_lighting_test_pb.GetReportCSVRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_lighting_test_pb.ReportCSV>;

}

