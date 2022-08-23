import {grpcWebEndpoint} from '@/api/config.js';
import Vue from 'vue';

export function closeResource(resource) {
  if (resource?.stream?.close) resource.stream.close();
}

export function setValue(resource, val) {
  Vue.set(resource, 'loading', false);
  Vue.set(resource, 'streamError', null);
  Vue.set(resource, 'value', val);
  Vue.set(resource, 'updateTime', new Date());
}

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
      Vue.delete(resource.value, idFunc(oldV))
    }
  }
  Vue.set(resource, 'updateTime', change.getChangeTime().toDate());
}

export function setError(resource, err) {
  Vue.set(resource, 'loading', false);
  Vue.set(resource, 'streamError', err);
  Vue.set(resource, 'updateTime', new Date());
}

export function pullResource(logPrefix, resource, newStream) {
  const doPull = (retryDelayMs = 1000) => {
    let retryCalled = false;
    const retry = () => {
      if (retryCalled) return;
      retryCalled = true;

      const handle = setTimeout(() => {
        const delay = Math.max(1000, Math.min(retryDelayMs * 2, 15 * 1000));
        doPull(delay);
      }, retryDelayMs)
      // fake stream we use to cancel the timeout if this component is disposed.
      Vue.set(resource, 'stream', {
        cancel() {
          clearTimeout(handle)
        }
      })
    };

    Promise.resolve(grpcWebEndpoint())
        .then(endpoint => {
          const stream = newStream(endpoint)
          Vue.set(resource, 'stream', stream);
          stream.on('data', () => {
            retryDelayMs = 1000; // if we were successful, we reset the retry delay
          });
          stream.on('error', err => {
            console.log(logPrefix, 'stream error', err);
            setError(resource, err);
            retry();
          });
          stream.on('end', () => {
            console.log(logPrefix, 'stream done');
            retry();
          })
        })
        .catch(err => {
          console.log(logPrefix, 'error caught', err);
          setError(resource, err);
          retry();
        });
  }

  doPull(0);
}

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

export function newActionTracker() {
  return {};
}

export function newResourceValue() {
  return {};
}

export function newResourceCollection() {
  return {};
}
