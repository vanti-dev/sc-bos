import * as grpcWeb from 'grpc-web';

import * as point_pb from './point_pb';


export class PointApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getPoints(
    request: point_pb.GetPointsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: point_pb.Points) => void
  ): grpcWeb.ClientReadableStream<point_pb.Points>;

  pullPoints(
    request: point_pb.PullPointsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<point_pb.PullPointsResponse>;

}

export class PointInfoClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describePoints(
    request: point_pb.DescribePointsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: point_pb.PointsSupport) => void
  ): grpcWeb.ClientReadableStream<point_pb.PointsSupport>;

}

export class PointApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getPoints(
    request: point_pb.GetPointsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<point_pb.Points>;

  pullPoints(
    request: point_pb.PullPointsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<point_pb.PullPointsResponse>;

}

export class PointInfoPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describePoints(
    request: point_pb.DescribePointsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<point_pb.PointsSupport>;

}

