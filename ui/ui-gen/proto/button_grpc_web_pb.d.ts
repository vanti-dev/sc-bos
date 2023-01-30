import * as grpcWeb from 'grpc-web';

import * as proto_button_pb from '../proto/button_pb';


export class ButtonApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getButtonState(
    request: proto_button_pb.GetButtonStateRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_button_pb.ButtonState) => void
  ): grpcWeb.ClientReadableStream<proto_button_pb.ButtonState>;

  pullButtonEvents(
    request: proto_button_pb.PullButtonEventsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_button_pb.PullButtonEventsResponse>;

}

export class ButtonApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getButtonState(
    request: proto_button_pb.GetButtonStateRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_button_pb.ButtonState>;

  pullButtonEvents(
    request: proto_button_pb.PullButtonEventsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_button_pb.PullButtonEventsResponse>;

}

