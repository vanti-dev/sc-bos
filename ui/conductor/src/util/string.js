
/**
 * Inserts spaces into camelCased string
 *
 * @param {string} key
 * @return {string}
 */
export function camelToSentence(key) {
  return key.replace(/([A-Z])/g, ' $1');
}
