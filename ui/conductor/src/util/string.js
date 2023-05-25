/**
 * Inserts spaces into camelCased string
 *
 * @param {string} key
 * @return {string}
 */
export function camelToSentence(key) {
  return key.replace(/([A-Z])/g, ' $1');
}

/**
 *
 * @param {string} string
 * @return {string}
 */
export function capitaliseString(string) {
  const firstChar = string.charAt(0).toUpperCase();
  const remainingChars = string.slice(1);
  const capitalisedString = firstChar + remainingChars;

  return capitalisedString;
}
