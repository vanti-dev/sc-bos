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

/**
 * Decides if a system ID is for the gateway system.
 *
 * @param {string} id
 * @return {boolean}
 */
export function isGatewayId(id) {
  return id === 'gateway' || id === 'proxy';
}
