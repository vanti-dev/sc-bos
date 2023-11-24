/**
 *
 * @param {string} message
 * @return {string}
 */
export function formatErrorMessage(message) {
  // Finding the start of the description and extracting it
  const descStart = message.indexOf('desc = ');
  const description = message.substring(descStart).split('"')[1];

  // Breaking down the description into parts
  const descParts = description.split(': ');
  // The last part is usually the main error detail
  const specificError = descParts[descParts.length - 1].trim();

  // Formatting the message into a more structured and readable sentence
  return `${specificError}.`;
}

