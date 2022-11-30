/**
 * Copy props from src into dst. Src is a JS object and dst is a protobuf message.
 * This function calls setSomeProp for each prop, which should be camelCase versions of the proto property names.
 *
 * @param {import('google-protobuf').Message} dst
 * @param {Object} src
 * @param {string} props
 */
import {Duration} from 'google-protobuf/google/protobuf/duration_pb.js';
import {FieldMask} from 'google-protobuf/google/protobuf/field_mask_pb.js';

export function setProperties(dst, src, ...props) {
  for (const prop of props) {
    if (src.hasOwnProperty(prop)) {
      dst['set' + prop[0].toUpperCase() + prop.substring(1)](src[prop]);
    }
  }
}

/**
 * Covert a JS object representing a Duration into a protobuf Duration.
 *
 * @param {Duration.AsObject} obj
 * @return {Duration}
 */
export function durationFromObject(obj) {
  if (!obj) return undefined;

  const dst = new Duration();
  setProperties(dst, obj, 'seconds', 'nanos');
  return dst;
}

/**
 * Covert a JS object representing a FieldMask into a protobuf FieldMask.
 *
 * @param {FieldMask.AsObject} obj
 * @return {FieldMask|undefined}
 */
export function fieldMaskFromObject(obj) {
  if (!obj) return undefined;

  const dst = new FieldMask();
  setProperties(dst, obj, 'pathsList');
  return dst;
}
