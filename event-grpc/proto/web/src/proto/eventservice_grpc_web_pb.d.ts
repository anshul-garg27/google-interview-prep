import * as grpcWeb from 'grpc-web';

import * as proto_eventservice_pb from '../proto/eventservice_pb';


export class EventServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  dispatch(
    request: proto_eventservice_pb.Events,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: proto_eventservice_pb.Response) => void
  ): grpcWeb.ClientReadableStream<proto_eventservice_pb.Response>;

}

export class EventServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  dispatch(
    request: proto_eventservice_pb.Events,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_eventservice_pb.Response>;

}

