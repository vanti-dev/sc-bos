import * as grpcWeb from 'grpc-web';

import * as proto_alerts_pb from '../proto/alerts_pb';


export class AlertApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listAlerts(
    request: proto_alerts_pb.ListAlertsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_alerts_pb.ListAlertsResponse) => void
  ): grpcWeb.ClientReadableStream<proto_alerts_pb.ListAlertsResponse>;

  pullAlerts(
    request: proto_alerts_pb.PullAlertsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_alerts_pb.PullAlertsResponse>;

  acknowledgeAlert(
    request: proto_alerts_pb.AcknowledgeAlertRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_alerts_pb.Alert) => void
  ): grpcWeb.ClientReadableStream<proto_alerts_pb.Alert>;

  unacknowledgeAlert(
    request: proto_alerts_pb.AcknowledgeAlertRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_alerts_pb.Alert) => void
  ): grpcWeb.ClientReadableStream<proto_alerts_pb.Alert>;

}

export class AlertAdminApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  createAlert(
    request: proto_alerts_pb.CreateAlertRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_alerts_pb.Alert) => void
  ): grpcWeb.ClientReadableStream<proto_alerts_pb.Alert>;

  updateAlert(
    request: proto_alerts_pb.UpdateAlertRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_alerts_pb.Alert) => void
  ): grpcWeb.ClientReadableStream<proto_alerts_pb.Alert>;

  deleteAlert(
    request: proto_alerts_pb.DeleteAlertRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_alerts_pb.DeleteAlertResponse) => void
  ): grpcWeb.ClientReadableStream<proto_alerts_pb.DeleteAlertResponse>;

}

export class AlertApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listAlerts(
    request: proto_alerts_pb.ListAlertsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_alerts_pb.ListAlertsResponse>;

  pullAlerts(
    request: proto_alerts_pb.PullAlertsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_alerts_pb.PullAlertsResponse>;

  acknowledgeAlert(
    request: proto_alerts_pb.AcknowledgeAlertRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_alerts_pb.Alert>;

  unacknowledgeAlert(
    request: proto_alerts_pb.AcknowledgeAlertRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_alerts_pb.Alert>;

}

export class AlertAdminApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  createAlert(
    request: proto_alerts_pb.CreateAlertRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_alerts_pb.Alert>;

  updateAlert(
    request: proto_alerts_pb.UpdateAlertRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_alerts_pb.Alert>;

  deleteAlert(
    request: proto_alerts_pb.DeleteAlertRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_alerts_pb.DeleteAlertResponse>;

}

