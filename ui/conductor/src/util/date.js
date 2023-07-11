
/**
 * Truncates the given date to the nearest round.
 * A date of 12:34 rounded down with 15minutes will return 12:30.
 *
 * @param {Date} date
 * @param {number} round
 * @return {Date}
 */
export function roundDown(date, round) {
  return new Date(date.getTime() - (date.getTime() % round));
}

/**
 * Truncates the given date to the nearest round.
 * A date of 12:34 rounded up with 15minutes will return 12:45.
 *
 * @param {Date} date
 * @param {number} round
 * @return {Date}
 */
export function roundUp(date, round) {
  const d = new Date(date.getTime() - (date.getTime() % round));
  if (d < date) {
    return new Date(d.getTime() + round);
  }
  return d;
}
