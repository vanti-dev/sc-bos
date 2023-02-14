import * as grpcWeb from 'grpc-web';

import * as proto_udmi_pb from '../proto/udmi_pb';


export class UdmiServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  pullControlTopics(
    request: proto_udmi_pb.PullControlTopicsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_udmi_pb.PullControlTopicsResponse>;

  onMessage(
    request: proto_udmi_pb.OnMessageRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_udmi_pb.OnMessageResponse) => void
  ): grpcWeb.ClientReadableStream<proto_udmi_pb.OnMessageResponse>;

  pullExportMessages(
    request: proto_udmi_pb.PullExportMessagesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_udmi_pb.PullExportMessagesResponse>;

}

export class UdmiServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  pullControlTopics(
    request: proto_udmi_pb.PullControlTopicsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_udmi_pb.PullControlTopicsResponse>;

  onMessage(
    request: proto_udmi_pb.OnMessageRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_udmi_pb.OnMessageResponse>;

  pullExportMessages(
    request: proto_udmi_pb.PullExportMessagesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_udmi_pb.PullExportMessagesResponse>;

}

