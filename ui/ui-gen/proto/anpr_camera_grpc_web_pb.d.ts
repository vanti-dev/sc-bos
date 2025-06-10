import * as grpcWeb from 'grpc-web';

import * as anpr_camera_pb from './anpr_camera_pb'; // proto import: "anpr_camera.proto"


export class AnprCameraApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listAnprEvents(
    request: anpr_camera_pb.ListAnprEventsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: anpr_camera_pb.ListAnprEventsResponse) => void
  ): grpcWeb.ClientReadableStream<anpr_camera_pb.ListAnprEventsResponse>;

  pullAnprEvents(
    request: anpr_camera_pb.PullAnprEventsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<anpr_camera_pb.PullAnprEventsResponse>;

}

export class AnprCameraApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listAnprEvents(
    request: anpr_camera_pb.ListAnprEventsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<anpr_camera_pb.ListAnprEventsResponse>;

  pullAnprEvents(
    request: anpr_camera_pb.PullAnprEventsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<anpr_camera_pb.PullAnprEventsResponse>;

}

