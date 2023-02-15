import * as grpcWeb from 'grpc-web';

import * as alerts_pb from './alerts_pb';


export class AlertApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listAlerts(
    request: alerts_pb.ListAlertsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: alerts_pb.ListAlertsResponse) => void
  ): grpcWeb.ClientReadableStream<alerts_pb.ListAlertsResponse>;

  pullAlerts(
    request: alerts_pb.PullAlertsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<alerts_pb.PullAlertsResponse>;

  acknowledgeAlert(
    request: alerts_pb.AcknowledgeAlertRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: alerts_pb.Alert) => void
  ): grpcWeb.ClientReadableStream<alerts_pb.Alert>;

  unacknowledgeAlert(
    request: alerts_pb.AcknowledgeAlertRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: alerts_pb.Alert) => void
  ): grpcWeb.ClientReadableStream<alerts_pb.Alert>;

  getAlertMetadata(
    request: alerts_pb.GetAlertMetadataRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: alerts_pb.AlertMetadata) => void
  ): grpcWeb.ClientReadableStream<alerts_pb.AlertMetadata>;

  pullAlertMetadata(
    request: alerts_pb.PullAlertMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<alerts_pb.PullAlertMetadataResponse>;

}

export class AlertAdminApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  createAlert(
    request: alerts_pb.CreateAlertRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: alerts_pb.Alert) => void
  ): grpcWeb.ClientReadableStream<alerts_pb.Alert>;

  updateAlert(
    request: alerts_pb.UpdateAlertRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: alerts_pb.Alert) => void
  ): grpcWeb.ClientReadableStream<alerts_pb.Alert>;

  deleteAlert(
    request: alerts_pb.DeleteAlertRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: alerts_pb.DeleteAlertResponse) => void
  ): grpcWeb.ClientReadableStream<alerts_pb.DeleteAlertResponse>;

}

export class AlertApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listAlerts(
    request: alerts_pb.ListAlertsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<alerts_pb.ListAlertsResponse>;

  pullAlerts(
    request: alerts_pb.PullAlertsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<alerts_pb.PullAlertsResponse>;

  acknowledgeAlert(
    request: alerts_pb.AcknowledgeAlertRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<alerts_pb.Alert>;

  unacknowledgeAlert(
    request: alerts_pb.AcknowledgeAlertRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<alerts_pb.Alert>;

  getAlertMetadata(
    request: alerts_pb.GetAlertMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<alerts_pb.AlertMetadata>;

  pullAlertMetadata(
    request: alerts_pb.PullAlertMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<alerts_pb.PullAlertMetadataResponse>;

}

export class AlertAdminApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  createAlert(
    request: alerts_pb.CreateAlertRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<alerts_pb.Alert>;

  updateAlert(
    request: alerts_pb.UpdateAlertRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<alerts_pb.Alert>;

  deleteAlert(
    request: alerts_pb.DeleteAlertRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<alerts_pb.DeleteAlertResponse>;

}

