import * as grpcWeb from 'grpc-web';

import * as proto_devices_pb from '../proto/devices_pb';


export class DevicesApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listDevices(
    request: proto_devices_pb.ListDevicesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: proto_devices_pb.ListDevicesResponse) => void
  ): grpcWeb.ClientReadableStream<proto_devices_pb.ListDevicesResponse>;

  pullDevices(
    request: proto_devices_pb.PullDevicesRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<proto_devices_pb.PullDevicesResponse>;

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

}

