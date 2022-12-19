import {grpcWebEndpoint} from '@/api/config.js';
import Vue from 'vue';

/**
 *
 * @param {RemoteResource} resource
 */
export function closeResource(resource) {
  if (resource?.stream?.close) resource.stream.close();
}

/**
 *
 * @param {ResourceValue} resource
 * @param {*} val
 */
export function setValue(resource, val) {
  Vue.set(resource, 'loading', false);
  Vue.set(resource, 'streamError', null);
  Vue.set(resource, 'value', val);
  Vue.set(resource, 'updateTime', new Date());
}

/**
 * @function IdFunc
 * @param {*} val
 * @return {string|number}
 */

/**
 *
 * @param {ResourceValue} resource
 * @param {Change} change
 * @param {IdFunc} idFunc
 */
export function setCollection(resource, change, idFunc) {
  Vue.set(resource, 'loading', false);
  Vue.set(resource, 'streamError', null);
  const oldV = change.getOldValue()?.toObject();
  const newV = change.getNewValue()?.toObject();
  if (newV) {
    if (!resource.value) Vue.set(resource, 'value', {});
    Vue.set(resource.value, idFunc(newV), newV);
  } else if (oldV) {
    if (resource.value) {
      Vue.delete(resource.value, idFunc(oldV));
    }
  }
  Vue.set(resource, 'updateTime', change.getChangeTime().toDate());
}

/**
 *
 * @param {RemoteResource} resource
 * @param {Error} err
 */
export function setError(resource, err) {
  Vue.set(resource, 'loading', false);
  Vue.set(resource, 'streamError', err);
  Vue.set(resource, 'updateTime', new Date());
}

/**
 *
 * @param {string} logPrefix
 * @param {RemoteResource} resource
 * @param {StreamFactory} newStream
 */
export function pullResource(logPrefix, resource, newStream) {
  const doPull = (retryDelayMs = 1000) => {
    let retryCalled = false;
    const retry = () => {
      if (retryCalled) return;
      retryCalled = true;

      const handle = setTimeout(() => {
        const delay = Math.max(1000, Math.min(retryDelayMs * 2, 15 * 1000));
        doPull(delay);
      }, retryDelayMs);
      // fake stream we use to cancel the timeout if this component is disposed.
      Vue.set(resource, 'stream', {
        cancel() {
          clearTimeout(handle);
        }
      });
    };

    Promise.resolve(grpcWebEndpoint())
        .then((endpoint) => {
          const stream = newStream(endpoint);
          Vue.set(resource, 'stream', stream);
          stream.on('data', () => {
            retryDelayMs = 1000; // if we were successful, we reset the retry delay
          });
          stream.on('error', (err) => {
            setError(resource, err);
            retry();
          });
          stream.on('end', () => {
            retry();
          });
        })
        .catch((err) => {
          setError(resource, err);
          retry();
        });
  };

  doPull(0);
}

/**
 *
 * @param {string} logPrefix
 * @param {ActionTracker} tracker
 * @param {Action} action
 */
export async function trackAction(logPrefix, tracker, action) {
  Vue.set(tracker, 'loading', true);
  const endpoint = await grpcWebEndpoint();
  try {
    const msg = await action(endpoint);
    const value = msg.toObject();
    Vue.set(tracker, 'response', value);
    return value;
  } catch (err) {
    Vue.set(tracker, 'error', err);
    throw err;
  } finally {
    Vue.set(tracker, 'loading', false);
  }
}

/**
 *
 * @return {ActionTracker}
 */
export function newActionTracker() {
  return {
    loading: false,
    response: null,
    error: null,
    duration: 0
  };
}

/**
 *
 * @return {ResourceValue}
 */
export function newResourceValue() {
  return {
    loading: false,
    stream: null,
    streamError: null,
    updateTime: null,
    value: null
  };
}

/**
 *
 * @return {ResourceCollection}
 */
export function newResourceCollection() {
  return {
    loading: false,
    stream: null,
    streamError: null,
    updateTime: null,
    value: {}
  };
}
