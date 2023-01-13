import {clientOptions} from '@/api/grpcweb.js';
import {pullResource, setValue} from '@/api/resource.js';
import {AirTemperatureApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/air_temperature_grpc_web_pb';
import {AirTemperature, PullAirTemperatureRequest} from '@smart-core-os/sc-api-grpc-web/traits/air_temperature_pb';

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
    return stream;
  });
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
 * Prepare an AirTemperature object for display to the user
 *
 * @param {AirTemperature.AsObject} temperatureValue
 * @return {Object}
 */
export function toDisplayObject(temperatureValue) {
  const data = {};
  Object.entries(temperatureValue).forEach(([key, value]) => {
    if (value !== undefined) {
      switch (key) {
        case 'mode': {
          data[key] = airTemperatureModeToString(value);
          break;
        }
        case 'ambientTemperature': {
          data['currentTemp'] = temperatureToString(value);
          break;
        }
        case 'temperatureSetPoint': {
          data['setPoint'] = temperatureToString(value);
          break;
        }
        case 'ambientHumidity': {
          data['humidity'] = (value * 100).toFixed(1) + '%';
          break;
        }
        case 'dewPoint': {
          data[key] = temperatureToString(value);
          break;
        }
        default: {
          if (key.toLowerCase().startsWith('ambient')) {
            key = key.substring(7);
          }
          data[key] = value;
        }
      }
    }
  });
  return data;
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
