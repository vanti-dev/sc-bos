import {fieldMaskFromObject, setProperties} from '@/api/convpb.js';
import {clientOptions} from '@/api/grpcweb.js';
import {pullResource, setValue, trackAction} from '@/api/resource.js';
import {tweenFromObject} from '@/api/sc/types/tween.js';
import {LightApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/light_grpc_web_pb';
import {
  Brightness,
  GetBrightnessRequest,
  LightPreset,
  PullBrightnessRequest,
  UpdateBrightnessRequest
} from '@smart-core-os/sc-api-grpc-web/traits/light_pb';

/**
 * @param {string} name
 * @param {ResourceValue<Brightness.AsObject, Brightness>} resource
 */
export function pullBrightness(name, resource) {
  pullResource('Light.Brightness', resource, endpoint => {
    const api = new LightApiPromiseClient(endpoint, null, clientOptions());
    const stream = api.pullBrightness(new PullBrightnessRequest().setName(name));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getBrightness().toObject());
      }
    });
    return stream;
  });
}

/**
 * @param {GetBrightnessRequest.AsObject} request
 * @param {ActionTracker<Brightness.AsObject>} tracker
 * @return {Promise<Brightness.AsObject>}
 */
export function getBrightness(request, tracker) {
  return trackAction('Light.getBrightness', tracker ?? {}, endpoint => {
    const api = new LightApiPromiseClient(endpoint, null, clientOptions());
    return api.getBrightness(getBrightnessRequestFromObject(request));
  });
}

/**
 * @param {UpdateBrightnessRequest.AsObject} request
 * @param {ActionTracker<Brightness.AsObject>} tracker
 * @return {Promise<Brightness.AsObject>}
 */
export function updateBrightness(request, tracker) {
  return trackAction('Light.updateBrightness', tracker ?? {}, endpoint => {
    const api = new LightApiPromiseClient(endpoint, null, clientOptions());
    return api.updateBrightness(updateBrightnessRequestFromObject(request));
  });
}

/**
 *
 * @param {GetBrightnessRequest.AsObject} obj
 * @return {GetBrightnessRequest}
 */
export function getBrightnessRequestFromObject(obj) {
  if (!obj) return undefined;

  const dst = new GetBrightnessRequest();
  setProperties(dst, obj, 'name');
  dst.setUpdateMask(fieldMaskFromObject(obj.readMask));
  return dst;
}

/**
 *
 * @param {UpdateBrightnessRequest.AsObject} obj
 * @return {UpdateBrightnessRequest}
 */
export function updateBrightnessRequestFromObject(obj) {
  if (!obj) return undefined;

  const dst = new UpdateBrightnessRequest();
  setProperties(dst, obj, 'name', 'delta');
  dst.setBrightness(brightnessFromObject(obj.brightness));
  dst.setUpdateMask(fieldMaskFromObject(obj.updateMask));
  return dst;
}

/**
 * Convert a JS object representation of Brightness into a protobuf Brightness object.
 *
 * @param {Brightness.AsObject} obj
 * @return {Brightness}
 */
export function brightnessFromObject(obj) {
  if (!obj) return undefined;

  const brightness = new Brightness();
  setProperties(brightness, obj, 'levelPercent', 'targetLevelPercent');
  brightness.setPreset(lightPresetFromObject(obj.preset));
  brightness.setTargetPreset(lightPresetFromObject(obj.targetPreset));
  brightness.setBrightnessTween(tweenFromObject(obj.brightnessTween));
  return brightness;
}

/**
 * Convert a JS object representation of LightPreset into a protobuf LightPreset object.
 *
 * @param {LightPreset.AsObject} obj
 * @return {LightPreset|null}
 */
export function lightPresetFromObject(obj) {
  if (!obj) return undefined;

  const dst = new LightPreset();
  setProperties(dst, obj, 'name', 'title');
  return dst;
}
