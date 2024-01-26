import * as grpcWeb from 'grpc-web';

import * as enrollment_pb from './enrollment_pb'; // proto import: "enrollment.proto"


export class EnrollmentApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getEnrollment(
    request: enrollment_pb.GetEnrollmentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: enrollment_pb.Enrollment) => void
  ): grpcWeb.ClientReadableStream<enrollment_pb.Enrollment>;

  createEnrollment(
    request: enrollment_pb.CreateEnrollmentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: enrollment_pb.Enrollment) => void
  ): grpcWeb.ClientReadableStream<enrollment_pb.Enrollment>;

  updateEnrollment(
    request: enrollment_pb.UpdateEnrollmentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: enrollment_pb.Enrollment) => void
  ): grpcWeb.ClientReadableStream<enrollment_pb.Enrollment>;

  deleteEnrollment(
    request: enrollment_pb.DeleteEnrollmentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: enrollment_pb.Enrollment) => void
  ): grpcWeb.ClientReadableStream<enrollment_pb.Enrollment>;

  testEnrollment(
    request: enrollment_pb.TestEnrollmentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: enrollment_pb.TestEnrollmentResponse) => void
  ): grpcWeb.ClientReadableStream<enrollment_pb.TestEnrollmentResponse>;

}

export class EnrollmentApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getEnrollment(
    request: enrollment_pb.GetEnrollmentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<enrollment_pb.Enrollment>;

  createEnrollment(
    request: enrollment_pb.CreateEnrollmentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<enrollment_pb.Enrollment>;

  updateEnrollment(
    request: enrollment_pb.UpdateEnrollmentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<enrollment_pb.Enrollment>;

  deleteEnrollment(
    request: enrollment_pb.DeleteEnrollmentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<enrollment_pb.Enrollment>;

  testEnrollment(
    request: enrollment_pb.TestEnrollmentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<enrollment_pb.TestEnrollmentResponse>;

}

