import * as grpcWeb from 'grpc-web';

import * as sound_sensor_pb from './sound_sensor_pb'; // proto import: "sound_sensor.proto"


export class SoundSensorApiClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getSoundLevel(
    request: sound_sensor_pb.GetSoundLevelRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: sound_sensor_pb.SoundLevel) => void
  ): grpcWeb.ClientReadableStream<sound_sensor_pb.SoundLevel>;

  pullSoundLevel(
    request: sound_sensor_pb.PullSoundLevelRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<sound_sensor_pb.PullSoundLevelResponse>;

}

export class SoundSensorInfoClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describeSoundLevel(
    request: sound_sensor_pb.DescribeSoundLevelRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: sound_sensor_pb.SoundLevelSupport) => void
  ): grpcWeb.ClientReadableStream<sound_sensor_pb.SoundLevelSupport>;

}

export class SoundSensorApiPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getSoundLevel(
    request: sound_sensor_pb.GetSoundLevelRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<sound_sensor_pb.SoundLevel>;

  pullSoundLevel(
    request: sound_sensor_pb.PullSoundLevelRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<sound_sensor_pb.PullSoundLevelResponse>;

}

export class SoundSensorInfoPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  describeSoundLevel(
    request: sound_sensor_pb.DescribeSoundLevelRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<sound_sensor_pb.SoundLevelSupport>;

}

