/**
 * Returns the directory containing the given path.
 *
 * @example
 * dir('a/b/c') // 'a/b'
 *
 * @param {string} path
 * @return {string}
 */
export function dir(path) {
  const i = path.lastIndexOf('/');
  return i === -1 ? path : path.slice(0, i);
}

/**
 * Returns dir(base)+path if path starts with './', otherwise returns path.
 *
 * @example
 * subPath('./b/c', 'a/foo') // 'a/b/c'
 *
 * @param {string} path
 * @param {string=} base
 * @return {string}
 */
export function subPath(path, base = '') {
  if (base !== '' && path?.startsWith('./')) {
    return dir(base) + path.slice(1);
  }
  return path;
}
