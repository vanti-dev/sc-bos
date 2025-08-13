import * as grpcWeb from 'grpc-web';

import * as report_pb from './report_pb'; // proto import: "report.proto"


export class ReportApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listReports(
    request: report_pb.ListReportsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: report_pb.ListReportsResponse) => void
  ): grpcWeb.ClientReadableStream<report_pb.ListReportsResponse>;

  getDownloadReportUrl(
    request: report_pb.GetDownloadReportUrlRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: report_pb.DownloadReportUrl) => void
  ): grpcWeb.ClientReadableStream<report_pb.DownloadReportUrl>;

}

export class ReportApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listReports(
    request: report_pb.ListReportsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<report_pb.ListReportsResponse>;

  getDownloadReportUrl(
    request: report_pb.GetDownloadReportUrlRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<report_pb.DownloadReportUrl>;

}

