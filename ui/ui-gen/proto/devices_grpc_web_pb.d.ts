import * as grpcWeb from 'grpc-web';

import * as proto_devices_pb from '../proto/devices_pb';


export class DevicesApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listDevices(
    request: proto_devices_pb.ListDevicesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_devices_pb.ListDevicesResponse) => void
  ): grpcWeb.ClientReadableStream<proto_devices_pb.ListDevicesResponse>;

  pullDevices(
    request: proto_devices_pb.PullDevicesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_devices_pb.PullDevicesResponse>;

  getDevicesMetadata(
    request: proto_devices_pb.GetDevicesMetadataRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_devices_pb.DevicesMetadata) => void
  ): grpcWeb.ClientReadableStream<proto_devices_pb.DevicesMetadata>;

  pullDevicesMetadata(
    request: proto_devices_pb.PullDevicesMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_devices_pb.PullDevicesMetadataResponse>;

}

export class DevicesApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listDevices(
    request: proto_devices_pb.ListDevicesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_devices_pb.ListDevicesResponse>;

  pullDevices(
    request: proto_devices_pb.PullDevicesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_devices_pb.PullDevicesResponse>;

  getDevicesMetadata(
    request: proto_devices_pb.GetDevicesMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_devices_pb.DevicesMetadata>;

  pullDevicesMetadata(
    request: proto_devices_pb.PullDevicesMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_devices_pb.PullDevicesMetadataResponse>;

}

