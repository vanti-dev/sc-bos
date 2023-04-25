/**
 *
 * @param {Promise<string>} controllerName
 * @param {string} serviceName
 * @return {string}
 */
export async function serviceName(controllerName, serviceName) {
  if (await controllerName === '') {
    return serviceName;
  } else {
    return await controllerName+ '/' + serviceName;
  }
}
