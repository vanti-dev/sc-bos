import * as grpcWeb from 'grpc-web';

import * as button_pb from './button_pb';


export class ButtonApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getButtonState(
    request: button_pb.GetButtonStateRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: button_pb.ButtonState) => void
  ): grpcWeb.ClientReadableStream<button_pb.ButtonState>;

  pullButtonState(
    request: button_pb.PullButtonStateRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<button_pb.PullButtonStateResponse>;

  updateButtonState(
    request: button_pb.UpdateButtonStateRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: button_pb.ButtonState) => void
  ): grpcWeb.ClientReadableStream<button_pb.ButtonState>;

}

export class ButtonApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getButtonState(
    request: button_pb.GetButtonStateRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<button_pb.ButtonState>;

  pullButtonState(
    request: button_pb.PullButtonStateRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<button_pb.PullButtonStateResponse>;

  updateButtonState(
    request: button_pb.UpdateButtonStateRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<button_pb.ButtonState>;

}

