/**
 *
 * @param {string} controllerName
 * @param {string} serviceName
 * @return {string}
 */
export function serviceName(controllerName, serviceName) {
  if (controllerName === '') {
    return controllerName;
  } else {
    return controllerName+ '/' + serviceName;
  }
}
