import * as grpcWeb from 'grpc-web';

import * as devices_pb from './devices_pb'; // proto import: "devices.proto"


export class DevicesApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listDevices(
    request: devices_pb.ListDevicesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: devices_pb.ListDevicesResponse) => void
  ): grpcWeb.ClientReadableStream<devices_pb.ListDevicesResponse>;

  pullDevices(
    request: devices_pb.PullDevicesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<devices_pb.PullDevicesResponse>;

  getDevicesMetadata(
    request: devices_pb.GetDevicesMetadataRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: devices_pb.DevicesMetadata) => void
  ): grpcWeb.ClientReadableStream<devices_pb.DevicesMetadata>;

  pullDevicesMetadata(
    request: devices_pb.PullDevicesMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<devices_pb.PullDevicesMetadataResponse>;

}

export class DevicesApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listDevices(
    request: devices_pb.ListDevicesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<devices_pb.ListDevicesResponse>;

  pullDevices(
    request: devices_pb.PullDevicesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<devices_pb.PullDevicesResponse>;

  getDevicesMetadata(
    request: devices_pb.GetDevicesMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<devices_pb.DevicesMetadata>;

  pullDevicesMetadata(
    request: devices_pb.PullDevicesMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<devices_pb.PullDevicesMetadataResponse>;

}

