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

/**
 *
 * @param {string} string
 * @return {string}
 */
export function camelCasingString(string) {
  let camelCasedString;

  if (string) camelCasedString = string.replace(/[^a-zA-Z0-9]+(.)/g, (_, chr) => chr.toUpperCase());

  return camelCasedString;
}
