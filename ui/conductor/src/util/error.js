/**
 *
 * @param {string} message
 * @return {string}
 */
export function formatErrorMessage(message) {
  // Ensure that message is a string
  if (typeof message !== 'string') {
    // Handle non-string message appropriately
    // For example, return a default error message or convert message to a string if possible
    return 'Invalid error message format.';
  }

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


