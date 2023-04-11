import {clientOptions} from '@/api/grpcweb.js';
import {pullResource, setError, setValue} from '@/api/resource.js';
import {AirTemperatureApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/air_temperature_grpc_web_pb';
import {
  AirTemperature,
  PullAirTemperatureRequest,
  UpdateAirTemperatureRequest
} from '@smart-core-os/sc-api-grpc-web/traits/air_temperature_pb';
import {trackAction} from '@/api/resource';
import {fieldMaskFromObject, setProperties} from '@/api/convpb';
import {Temperature} from '@smart-core-os/sc-api-grpc-web/types/unit_pb';

/**
 * @param {string} name
 * @param {ResourceValue<AirTemperature.AsObject, AirTemperature>} resource
 */
export function pullAirTemperature(name, resource) {
  pullResource('AirTemperature', resource, endpoint => {
    const api = new AirTemperatureApiPromiseClient(endpoint, null, clientOptions());
    const stream = api.pullAirTemperature(new PullAirTemperatureRequest().setName(name));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getAirTemperature().toObject());
      }
    });
    stream.on('error', err => {
      setError(resource, err);
    });
    return stream;
  });
}

/**
 *
 * @param {UpdateAirTemperatureRequest.AsObject} request
 * @param {ActionTracker<AirTemperature.AsObject>} tracker
 * @return {Promise<AirTemperature.AsObject>}
 */
export function updateAirTemperature(request, tracker) {
  return trackAction('AirTemperature.updateAirTemperature', tracker ?? {}, endpoint => {
    const api = new AirTemperatureApiPromiseClient(endpoint, null, clientOptions());
    return api.updateAirTemperature(updateAirTemperatureRequestFromObject(request));
  });
}

/**
 * @param {UpdateAirTemperatureRequest.AsObject} obj
 * @return {UpdateAirTemperatureRequest}
 */
export function updateAirTemperatureRequestFromObject(obj) {
  if (!obj) return undefined;

  const req = new UpdateAirTemperatureRequest();
  setProperties(req, obj, 'name');
  req.setState(stateFromObject(obj.state));
  req.setUpdateMask(fieldMaskFromObject(obj.updateMask));
  return req;
}

/**
 * @param {AirTemperature.AsObject} obj
 * @return {AirTemperature}
 */
export function stateFromObject(obj) {
  if (!obj) return undefined;

  const state = new AirTemperature();
  setProperties(state, obj, 'ambient_humidity');
  state.setTemperatureSetPoint(temperatureFromObject(obj.temperatureSetPoint));
  return state;
}

/**
 * @param {Temperature.AsObject} obj
 * @return {Temperature}
 */
export function temperatureFromObject(obj) {
  if (!obj) return undefined;

  const t = new Temperature();
  setProperties(t, obj, 'valueCelsius');
  return t;
}

/**
 * @param {AirTemperature.Mode} mode
 * @return {string}
 */
export function airTemperatureModeToString(mode) {
  switch (mode) {
    case AirTemperature.Mode.MODE_UNSPECIFIED: return 'Unknown';
    case AirTemperature.Mode.ON: return 'On';
    case AirTemperature.Mode.OFF: return 'Off';
    case AirTemperature.Mode.HEAT: return 'Heat';
    case AirTemperature.Mode.COOL: return 'Cool';
    case AirTemperature.Mode.HEAT_COOL: return 'Heat/Cool';
    case AirTemperature.Mode.AUTO: return 'Auto';
    case AirTemperature.Mode.FAN_ONLY: return 'Fan Only';
    case AirTemperature.Mode.ECO: return 'Eco';
    case AirTemperature.Mode.PURIFIER: return 'Purifier';
    case AirTemperature.Mode.DRY: return 'Dry';
    case AirTemperature.Mode.LOCKED: return 'Locked';
  }
}

/**
 *
 * @param {Temperature.AsObject} value
 * @return {string}
 */
export function temperatureToString(value) {
  if (value.hasOwnProperty('valueCelsius')) {
    return value.valueCelsius.toFixed(1) + '°C';
  }
  return '-';
}
