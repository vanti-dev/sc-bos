import * as grpcWeb from 'grpc-web';

import * as service_ticket_pb from './service_ticket_pb'; // proto import: "service_ticket.proto"


export class ServiceTicketApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  createTicket(
    request: service_ticket_pb.CreateTicketRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: service_ticket_pb.Ticket) => void
  ): grpcWeb.ClientReadableStream<service_ticket_pb.Ticket>;

  updateTicket(
    request: service_ticket_pb.UpdateTicketRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: service_ticket_pb.Ticket) => void
  ): grpcWeb.ClientReadableStream<service_ticket_pb.Ticket>;

}

export class ServiceTicketApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  createTicket(
    request: service_ticket_pb.CreateTicketRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<service_ticket_pb.Ticket>;

  updateTicket(
    request: service_ticket_pb.UpdateTicketRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<service_ticket_pb.Ticket>;

}

