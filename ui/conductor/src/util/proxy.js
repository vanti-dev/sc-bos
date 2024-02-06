/**
 *
 * @param {string} controllerName
 * @param {string} serviceName
 * @return {string}
 */
export function serviceName(controllerName, serviceName) {
  if (!controllerName || controllerName === '') {
    return serviceName;
  } else {
    return controllerName + '/' + serviceName;
  }
}
