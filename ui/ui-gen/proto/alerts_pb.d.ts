import * as jspb from 'google-protobuf'

import * as google_protobuf_field_mask_pb from 'google-protobuf/google/protobuf/field_mask_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as types_change_pb from '@smart-core-os/sc-api-grpc-web/types/change_pb';


export class Alert extends jspb.Message {
  getId(): string;
  setId(value: string): Alert;

  getDescription(): string;
  setDescription(value: string): Alert;

  getCreateTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreateTime(value?: google_protobuf_timestamp_pb.Timestamp): Alert;
  hasCreateTime(): boolean;
  clearCreateTime(): Alert;

  getAcknowledgement(): Alert.Acknowledgement | undefined;
  setAcknowledgement(value?: Alert.Acknowledgement): Alert;
  hasAcknowledgement(): boolean;
  clearAcknowledgement(): Alert;

  getSeverity(): Alert.Severity;
  setSeverity(value: Alert.Severity): Alert;

  getFloor(): string;
  setFloor(value: string): Alert;

  getZone(): string;
  setZone(value: string): Alert;

  getSource(): string;
  setSource(value: string): Alert;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Alert.AsObject;
  static toObject(includeInstance: boolean, msg: Alert): Alert.AsObject;
  static serializeBinaryToWriter(message: Alert, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Alert;
  static deserializeBinaryFromReader(message: Alert, reader: jspb.BinaryReader): Alert;
}

export namespace Alert {
  export type AsObject = {
    id: string,
    description: string,
    createTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    acknowledgement?: Alert.Acknowledgement.AsObject,
    severity: Alert.Severity,
    floor: string,
    zone: string,
    source: string,
  }

  export class Acknowledgement extends jspb.Message {
    getAcknowledgeTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setAcknowledgeTime(value?: google_protobuf_timestamp_pb.Timestamp): Acknowledgement;
    hasAcknowledgeTime(): boolean;
    clearAcknowledgeTime(): Acknowledgement;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Acknowledgement.AsObject;
    static toObject(includeInstance: boolean, msg: Acknowledgement): Acknowledgement.AsObject;
    static serializeBinaryToWriter(message: Acknowledgement, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Acknowledgement;
    static deserializeBinaryFromReader(message: Acknowledgement, reader: jspb.BinaryReader): Acknowledgement;
  }

  export namespace Acknowledgement {
    export type AsObject = {
      acknowledgeTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    }
  }


  export class Query extends jspb.Message {
    getCreatedNotBefore(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setCreatedNotBefore(value?: google_protobuf_timestamp_pb.Timestamp): Query;
    hasCreatedNotBefore(): boolean;
    clearCreatedNotBefore(): Query;

    getCreatedNotAfter(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setCreatedNotAfter(value?: google_protobuf_timestamp_pb.Timestamp): Query;
    hasCreatedNotAfter(): boolean;
    clearCreatedNotAfter(): Query;

    getSeverityNotBelow(): number;
    setSeverityNotBelow(value: number): Query;

    getSeverityNotAbove(): number;
    setSeverityNotAbove(value: number): Query;

    getFloor(): string;
    setFloor(value: string): Query;

    getZone(): string;
    setZone(value: string): Query;

    getSource(): string;
    setSource(value: string): Query;

    getAcknowledged(): boolean;
    setAcknowledged(value: boolean): Query;

    getAcknowledgedCase(): Query.AcknowledgedCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Query.AsObject;
    static toObject(includeInstance: boolean, msg: Query): Query.AsObject;
    static serializeBinaryToWriter(message: Query, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Query;
    static deserializeBinaryFromReader(message: Query, reader: jspb.BinaryReader): Query;
  }

  export namespace Query {
    export type AsObject = {
      createdNotBefore?: google_protobuf_timestamp_pb.Timestamp.AsObject,
      createdNotAfter?: google_protobuf_timestamp_pb.Timestamp.AsObject,
      severityNotBelow: number,
      severityNotAbove: number,
      floor: string,
      zone: string,
      source: string,
      acknowledged: boolean,
    }

    export enum AcknowledgedCase { 
      _ACKNOWLEDGED_NOT_SET = 0,
      ACKNOWLEDGED = 8,
    }
  }


  export enum Severity { 
    SEVERITY_UNSPECIFIED = 0,
    INFO = 9,
    WARNING = 13,
    SEVERE = 17,
    LIFE_SAFETY = 21,
  }
}

export class ListAlertsRequest extends jspb.Message {
  getName(): string;
  setName(value: string): ListAlertsRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): ListAlertsRequest;
  hasReadMask(): boolean;
  clearReadMask(): ListAlertsRequest;

  getPageSize(): number;
  setPageSize(value: number): ListAlertsRequest;

  getPageToken(): string;
  setPageToken(value: string): ListAlertsRequest;

  getQuery(): Alert.Query | undefined;
  setQuery(value?: Alert.Query): ListAlertsRequest;
  hasQuery(): boolean;
  clearQuery(): ListAlertsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAlertsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAlertsRequest): ListAlertsRequest.AsObject;
  static serializeBinaryToWriter(message: ListAlertsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAlertsRequest;
  static deserializeBinaryFromReader(message: ListAlertsRequest, reader: jspb.BinaryReader): ListAlertsRequest;
}

export namespace ListAlertsRequest {
  export type AsObject = {
    name: string,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
    pageSize: number,
    pageToken: string,
    query?: Alert.Query.AsObject,
  }
}

export class ListAlertsResponse extends jspb.Message {
  getAlertsList(): Array<Alert>;
  setAlertsList(value: Array<Alert>): ListAlertsResponse;
  clearAlertsList(): ListAlertsResponse;
  addAlerts(value?: Alert, index?: number): Alert;

  getNextPageToken(): string;
  setNextPageToken(value: string): ListAlertsResponse;

  getTotalSize(): number;
  setTotalSize(value: number): ListAlertsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAlertsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAlertsResponse): ListAlertsResponse.AsObject;
  static serializeBinaryToWriter(message: ListAlertsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAlertsResponse;
  static deserializeBinaryFromReader(message: ListAlertsResponse, reader: jspb.BinaryReader): ListAlertsResponse;
}

export namespace ListAlertsResponse {
  export type AsObject = {
    alertsList: Array<Alert.AsObject>,
    nextPageToken: string,
    totalSize: number,
  }
}

export class PullAlertsRequest extends jspb.Message {
  getName(): string;
  setName(value: string): PullAlertsRequest;

  getReadMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setReadMask(value?: google_protobuf_field_mask_pb.FieldMask): PullAlertsRequest;
  hasReadMask(): boolean;
  clearReadMask(): PullAlertsRequest;

  getQuery(): Alert.Query | undefined;
  setQuery(value?: Alert.Query): PullAlertsRequest;
  hasQuery(): boolean;
  clearQuery(): PullAlertsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullAlertsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PullAlertsRequest): PullAlertsRequest.AsObject;
  static serializeBinaryToWriter(message: PullAlertsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullAlertsRequest;
  static deserializeBinaryFromReader(message: PullAlertsRequest, reader: jspb.BinaryReader): PullAlertsRequest;
}

export namespace PullAlertsRequest {
  export type AsObject = {
    name: string,
    readMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
    query?: Alert.Query.AsObject,
  }
}

export class PullAlertsResponse extends jspb.Message {
  getChangesList(): Array<PullAlertsResponse.Change>;
  setChangesList(value: Array<PullAlertsResponse.Change>): PullAlertsResponse;
  clearChangesList(): PullAlertsResponse;
  addChanges(value?: PullAlertsResponse.Change, index?: number): PullAlertsResponse.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PullAlertsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PullAlertsResponse): PullAlertsResponse.AsObject;
  static serializeBinaryToWriter(message: PullAlertsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PullAlertsResponse;
  static deserializeBinaryFromReader(message: PullAlertsResponse, reader: jspb.BinaryReader): PullAlertsResponse;
}

export namespace PullAlertsResponse {
  export type AsObject = {
    changesList: Array<PullAlertsResponse.Change.AsObject>,
  }

  export class Change extends jspb.Message {
    getName(): string;
    setName(value: string): Change;

    getType(): types_change_pb.ChangeType;
    setType(value: types_change_pb.ChangeType): Change;

    getNewValue(): Alert | undefined;
    setNewValue(value?: Alert): Change;
    hasNewValue(): boolean;
    clearNewValue(): Change;

    getOldValue(): Alert | undefined;
    setOldValue(value?: Alert): Change;
    hasOldValue(): boolean;
    clearOldValue(): Change;

    getChangeTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setChangeTime(value?: google_protobuf_timestamp_pb.Timestamp): Change;
    hasChangeTime(): boolean;
    clearChangeTime(): Change;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Change.AsObject;
    static toObject(includeInstance: boolean, msg: Change): Change.AsObject;
    static serializeBinaryToWriter(message: Change, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Change;
    static deserializeBinaryFromReader(message: Change, reader: jspb.BinaryReader): Change;
  }

  export namespace Change {
    export type AsObject = {
      name: string,
      type: types_change_pb.ChangeType,
      newValue?: Alert.AsObject,
      oldValue?: Alert.AsObject,
      changeTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    }
  }

}

export class AcknowledgeAlertRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AcknowledgeAlertRequest;

  getId(): string;
  setId(value: string): AcknowledgeAlertRequest;

  getAllowAcknowledged(): boolean;
  setAllowAcknowledged(value: boolean): AcknowledgeAlertRequest;

  getAllowMissing(): boolean;
  setAllowMissing(value: boolean): AcknowledgeAlertRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AcknowledgeAlertRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AcknowledgeAlertRequest): AcknowledgeAlertRequest.AsObject;
  static serializeBinaryToWriter(message: AcknowledgeAlertRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AcknowledgeAlertRequest;
  static deserializeBinaryFromReader(message: AcknowledgeAlertRequest, reader: jspb.BinaryReader): AcknowledgeAlertRequest;
}

export namespace AcknowledgeAlertRequest {
  export type AsObject = {
    name: string,
    id: string,
    allowAcknowledged: boolean,
    allowMissing: boolean,
  }
}

export class CreateAlertRequest extends jspb.Message {
  getName(): string;
  setName(value: string): CreateAlertRequest;

  getAlert(): Alert | undefined;
  setAlert(value?: Alert): CreateAlertRequest;
  hasAlert(): boolean;
  clearAlert(): CreateAlertRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAlertRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAlertRequest): CreateAlertRequest.AsObject;
  static serializeBinaryToWriter(message: CreateAlertRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAlertRequest;
  static deserializeBinaryFromReader(message: CreateAlertRequest, reader: jspb.BinaryReader): CreateAlertRequest;
}

export namespace CreateAlertRequest {
  export type AsObject = {
    name: string,
    alert?: Alert.AsObject,
  }
}

export class UpdateAlertRequest extends jspb.Message {
  getName(): string;
  setName(value: string): UpdateAlertRequest;

  getAlert(): Alert | undefined;
  setAlert(value?: Alert): UpdateAlertRequest;
  hasAlert(): boolean;
  clearAlert(): UpdateAlertRequest;

  getUpdateMask(): google_protobuf_field_mask_pb.FieldMask | undefined;
  setUpdateMask(value?: google_protobuf_field_mask_pb.FieldMask): UpdateAlertRequest;
  hasUpdateMask(): boolean;
  clearUpdateMask(): UpdateAlertRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAlertRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAlertRequest): UpdateAlertRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateAlertRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAlertRequest;
  static deserializeBinaryFromReader(message: UpdateAlertRequest, reader: jspb.BinaryReader): UpdateAlertRequest;
}

export namespace UpdateAlertRequest {
  export type AsObject = {
    name: string,
    alert?: Alert.AsObject,
    updateMask?: google_protobuf_field_mask_pb.FieldMask.AsObject,
  }
}

export class DeleteAlertRequest extends jspb.Message {
  getName(): string;
  setName(value: string): DeleteAlertRequest;

  getId(): string;
  setId(value: string): DeleteAlertRequest;

  getAllowMissing(): boolean;
  setAllowMissing(value: boolean): DeleteAlertRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAlertRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAlertRequest): DeleteAlertRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteAlertRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAlertRequest;
  static deserializeBinaryFromReader(message: DeleteAlertRequest, reader: jspb.BinaryReader): DeleteAlertRequest;
}

export namespace DeleteAlertRequest {
  export type AsObject = {
    name: string,
    id: string,
    allowMissing: boolean,
  }
}

export class DeleteAlertResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAlertResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAlertResponse): DeleteAlertResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteAlertResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAlertResponse;
  static deserializeBinaryFromReader(message: DeleteAlertResponse, reader: jspb.BinaryReader): DeleteAlertResponse;
}

export namespace DeleteAlertResponse {
  export type AsObject = {
  }
}

