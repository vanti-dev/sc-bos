import * as grpcWeb from 'grpc-web';

import * as udmi_pb from './udmi_pb';


export class UdmiServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  pullControlTopics(
    request: udmi_pb.PullControlTopicsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<udmi_pb.PullControlTopicsResponse>;

  onMessage(
    request: udmi_pb.OnMessageRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: udmi_pb.OnMessageResponse) => void
  ): grpcWeb.ClientReadableStream<udmi_pb.OnMessageResponse>;

  pullExportMessages(
    request: udmi_pb.PullExportMessagesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<udmi_pb.PullExportMessagesResponse>;

  getExportMessage(
    request: udmi_pb.GetExportMessageRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: udmi_pb.MqttMessage) => void
  ): grpcWeb.ClientReadableStream<udmi_pb.MqttMessage>;

}

export class UdmiServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  pullControlTopics(
    request: udmi_pb.PullControlTopicsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<udmi_pb.PullControlTopicsResponse>;

  onMessage(
    request: udmi_pb.OnMessageRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<udmi_pb.OnMessageResponse>;

  pullExportMessages(
    request: udmi_pb.PullExportMessagesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<udmi_pb.PullExportMessagesResponse>;

  getExportMessage(
    request: udmi_pb.GetExportMessageRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<udmi_pb.MqttMessage>;

}

