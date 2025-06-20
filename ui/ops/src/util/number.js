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

/**
 * Scale a number from one range to another.
 *
 * @param {number} value
 * @param {number} fromMin
 * @param {number} fromMax
 * @param {number} toMin
 * @param {number} toMax
 * @return {number}
 */
export function scale(value, fromMin, fromMax, toMin, toMax) {
  return toMin + (toMax - toMin) * ((value - fromMin) / (fromMax - fromMin));
}

/**
 * Round a number to a specified number of decimal places.
 *
 * @param {number} num
 * @param {number} decimals
 * @return {number}
 */
export function roundTo(num, decimals) {
  if (decimals < 0) {
    return num;
  }

  if (decimals === 0) {
    return Math.round(num);
  }

  const factor = Math.pow(10, decimals);
  return Math.round(num * factor) / factor;
}