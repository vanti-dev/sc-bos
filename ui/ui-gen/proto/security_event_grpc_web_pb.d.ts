import * as grpcWeb from 'grpc-web';

import * as security_event_pb from './security_event_pb'; // proto import: "security_event.proto"


export class SecurityEventApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listSecurityEvents(
    request: security_event_pb.ListSecurityEventsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: security_event_pb.ListSecurityEventsResponse) => void
  ): grpcWeb.ClientReadableStream<security_event_pb.ListSecurityEventsResponse>;

  pullSecurityEvents(
    request: security_event_pb.PullSecurityEventsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<security_event_pb.PullSecurityEventsResponse>;

}

export class SecurityEventApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listSecurityEvents(
    request: security_event_pb.ListSecurityEventsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<security_event_pb.ListSecurityEventsResponse>;

  pullSecurityEvents(
    request: security_event_pb.PullSecurityEventsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<security_event_pb.PullSecurityEventsResponse>;

}

