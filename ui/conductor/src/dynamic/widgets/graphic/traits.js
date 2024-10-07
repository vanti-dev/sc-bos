/**
 * Attempts to extract the name from a query object.
 * Typically if the query isn't already a name then it will have a name property, so return that.
 *
 * @param {*} query
 * @return {undefined|string}
 */
export function nameFromRequest(query) {
  if (typeof query === 'string') return query;
  if (typeof query === 'object') return query.name;
  return undefined;
}
