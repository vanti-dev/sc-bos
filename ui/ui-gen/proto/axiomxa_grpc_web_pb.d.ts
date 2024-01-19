import * as grpcWeb from 'grpc-web';

import * as axiomxa_pb from './axiomxa_pb'; // proto import: "axiomxa.proto"


export class AxiomXaDriverServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  saveQRCredential(
    request: axiomxa_pb.SaveQRCredentialRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: axiomxa_pb.SaveQRCredentialResponse) => void
  ): grpcWeb.ClientReadableStream<axiomxa_pb.SaveQRCredentialResponse>;

}

export class AxiomXaDriverServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  saveQRCredential(
    request: axiomxa_pb.SaveQRCredentialRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<axiomxa_pb.SaveQRCredentialResponse>;

}

