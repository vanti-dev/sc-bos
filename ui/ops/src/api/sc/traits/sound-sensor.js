import {fieldMaskFromObject, setProperties} from '@/api/convpb.js';
import {pullResource, setValue, trackAction} from '@/api/resource.js';
import {SoundSensorApiPromiseClient, SoundSensorInfoPromiseClient} from '@vanti-dev/sc-bos-ui-gen/proto/sound_sensor_grpc_web_pb';
import {DescribeSoundLevelRequest, PullSoundLevelRequest} from '@vanti-dev/sc-bos-ui-gen/proto/sound_sensor_pb';

/**
 * @param {Partial<PullSoundLevelRequest.AsObject>} request
 * @param {ResourceValue<SoundLevel.AsObject, PullSoundLevelResponse>} resource
 */
export function pullSoundLevel(request, resource) {
 pullResource('SoundSensorApi.pullSoundLevel', resource, (endpoint) => {
   const api = apiClient(endpoint);
   const stream = api.pullSoundLevel(pullSoundLevelRequestFromObject(request));
   stream.on('data', (msg) => {
    const changes = msg.getChangesList();
    for (const change of changes) {
     setValue(resource, change.getSoundLevel().toObject());
    }
   });
   return stream;
  });
}

/**
 *
 * @param {Partial<DescribeSoundLevelRequest.AsObject>} request
 * @param {ActionTracker<SoundLevelSupport.AsObject>} [tracker]
 * @return {Promise<SoundLevelSupport.AsObject>}
 */
export function describeSoundLevel(request, tracker) {
  return trackAction('SoundLevelInfo.DescribeSoundLevel', tracker ?? {}, (endpoint) => {
   const api = infoClient(endpoint);
   return api.describeSoundLevel(describeSoundLevelRequestFromObject(request));
  });
}

/**
 * @param {string} endpoint
 * @return {SoundSensorApiPromiseClient}
 */
function apiClient(endpoint) {
  return new SoundSensorApiPromiseClient(endpoint, null, null);
}

/**
 * @param {string} endpoint
 * @return {SoundSensorInfoPromiseClient}
 */
function infoClient(endpoint) {
  return new SoundSensorInfoPromiseClient(endpoint, null, null);
}

/**
 * @param {Partial<PullSoundLevelRequest.AsObject>} obj
 * @return {PullSoundLevelRequest|undefined}
 */
function pullSoundLevelRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new PullSoundLevelRequest();
  setProperties(dst, obj, 'name', 'updatesOnly');
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  return dst;
}

/**
 * @param {Partial<DescribeSoundLevelRequest.AsObject>} obj
 * @return {undefined|DescribeSoundLevelRequest}
 */
function describeSoundLevelRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new DescribeSoundLevelRequest();
  setProperties(dst, obj, 'name');
  return dst;
}