import * as grpcWeb from 'grpc-web';

import * as proto_mqtt_pb from '../proto/mqtt_pb';


export class MqttServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  pullMessages(
    request: proto_mqtt_pb.PullMessagesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_mqtt_pb.PullMessagesResponse>;

}

export class MqttServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  pullMessages(
    request: proto_mqtt_pb.PullMessagesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_mqtt_pb.PullMessagesResponse>;

}

