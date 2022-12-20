import * as grpcWeb from 'grpc-web';

import * as proto_enrollment_pb from '../proto/enrollment_pb';


export class EnrollmentApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getEnrollment(
    request: proto_enrollment_pb.GetEnrollmentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_enrollment_pb.Enrollment) => void
  ): grpcWeb.ClientReadableStream<proto_enrollment_pb.Enrollment>;

  createEnrollment(
    request: proto_enrollment_pb.CreateEnrollmentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_enrollment_pb.Enrollment) => void
  ): grpcWeb.ClientReadableStream<proto_enrollment_pb.Enrollment>;

  deleteEnrollment(
    request: proto_enrollment_pb.DeleteEnrollmentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_enrollment_pb.Enrollment) => void
  ): grpcWeb.ClientReadableStream<proto_enrollment_pb.Enrollment>;

}

export class EnrollmentApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getEnrollment(
    request: proto_enrollment_pb.GetEnrollmentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_enrollment_pb.Enrollment>;

  createEnrollment(
    request: proto_enrollment_pb.CreateEnrollmentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_enrollment_pb.Enrollment>;

  deleteEnrollment(
    request: proto_enrollment_pb.DeleteEnrollmentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_enrollment_pb.Enrollment>;

}

