import * as jspb from 'google-protobuf'



export class Ticket extends jspb.Message {
  getId(): string;
  setId(value: string): Ticket;

  getDescription(): string;
  setDescription(value: string): Ticket;

  getDetails(): string;
  setDetails(value: string): Ticket;

  getReporterName(): string;
  setReporterName(value: string): Ticket;

  getClassification(): Ticket.Classification;
  setClassification(value: Ticket.Classification): Ticket;

  getSeverity(): Ticket.Severity;
  setSeverity(value: Ticket.Severity): Ticket;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Ticket.AsObject;
  static toObject(includeInstance: boolean, msg: Ticket): Ticket.AsObject;
  static serializeBinaryToWriter(message: Ticket, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Ticket;
  static deserializeBinaryFromReader(message: Ticket, reader: jspb.BinaryReader): Ticket;
}

export namespace Ticket {
  export type AsObject = {
    id: string,
    description: string,
    details: string,
    reporterName: string,
    classification: Ticket.Classification,
    severity: Ticket.Severity,
  }

  export enum Classification { 
    CLASSIFICATION_UNSPECIFIED = 0,
    MAINTENANCE = 1,
    CLEANING = 2,
    FIRE = 3,
    ELECTRICAL = 4,
    EMERGENCY_REPAIR = 5,
    NON_EMERGENCY_REPAIR = 6,
    PLUMBING = 7,
    WASTE = 8,
    RECYCLING = 9,
    SECURITY = 10,
    URGENT_REPAIR = 11,
    OBSERVATION = 12,
    PORTERAGE = 13,
    SPACE_PREPARATION = 14,
  }

  export enum Severity { 
    SEVERITY_UNSPECIFIED = 0,
    EMERGENCY = 1,
    CRITICAL = 2,
    URGENT = 3,
    HIGH = 4,
    MEDIUM = 5,
    LOW = 6,
  }
}

export class CreateTicketRequest extends jspb.Message {
  getName(): string;
  setName(value: string): CreateTicketRequest;

  getTicket(): Ticket | undefined;
  setTicket(value?: Ticket): CreateTicketRequest;
  hasTicket(): boolean;
  clearTicket(): CreateTicketRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateTicketRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateTicketRequest): CreateTicketRequest.AsObject;
  static serializeBinaryToWriter(message: CreateTicketRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateTicketRequest;
  static deserializeBinaryFromReader(message: CreateTicketRequest, reader: jspb.BinaryReader): CreateTicketRequest;
}

export namespace CreateTicketRequest {
  export type AsObject = {
    name: string,
    ticket?: Ticket.AsObject,
  }
}

export class UpdateTicketRequest extends jspb.Message {
  getName(): string;
  setName(value: string): UpdateTicketRequest;

  getTicket(): Ticket | undefined;
  setTicket(value?: Ticket): UpdateTicketRequest;
  hasTicket(): boolean;
  clearTicket(): UpdateTicketRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateTicketRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateTicketRequest): UpdateTicketRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateTicketRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateTicketRequest;
  static deserializeBinaryFromReader(message: UpdateTicketRequest, reader: jspb.BinaryReader): UpdateTicketRequest;
}

export namespace UpdateTicketRequest {
  export type AsObject = {
    name: string,
    ticket?: Ticket.AsObject,
  }
}

