import * as grpcWeb from 'grpc-web';

import * as anpr_camera_pb from './anpr_camera_pb'; // proto import: "anpr_camera.proto"


export class AnprCameraApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getEvent(
    request: anpr_camera_pb.GetLastEventRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: anpr_camera_pb.AnprEvent) => void
  ): grpcWeb.ClientReadableStream<anpr_camera_pb.AnprEvent>;

  pullEvents(
    request: anpr_camera_pb.PullEventsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<anpr_camera_pb.PullEventsResponse>;

}

export class AnprCameraApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getEvent(
    request: anpr_camera_pb.GetLastEventRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<anpr_camera_pb.AnprEvent>;

  pullEvents(
    request: anpr_camera_pb.PullEventsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<anpr_camera_pb.PullEventsResponse>;

}

