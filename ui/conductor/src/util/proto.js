/**
 * Converts a proto map, which is an array of [k,v] into a js object.
 *
 * @param {Array<[K,V]>} arr
 * @return {Object<K,V>}
 * @template K,V
 */
export function convertProtoMap(arr) {
  if (!arr) return {};
  const dst = {};
  for (const [k, v] of arr || []) {
    dst[k] = v;
  }
  return dst;
}
