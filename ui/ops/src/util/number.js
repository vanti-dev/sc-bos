export const MAX_INT32 = 2147483647; // 2^31 - 1

/**
 * Cap a number between min and max.
 *
 * @param {number} value
 * @param {number} min
 * @param {number} max
 * @return {number}
 */
export function cap(value, min, max) {
  return Math.max(min, Math.min(max, value));
}
