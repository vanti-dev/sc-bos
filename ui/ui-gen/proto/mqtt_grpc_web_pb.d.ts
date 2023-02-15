import * as grpcWeb from 'grpc-web';

import * as mqtt_pb from './mqtt_pb';


export class MqttServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  pullMessages(
    request: mqtt_pb.PullMessagesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<mqtt_pb.PullMessagesResponse>;

}

export class MqttServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  pullMessages(
    request: mqtt_pb.PullMessagesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<mqtt_pb.PullMessagesResponse>;

}

