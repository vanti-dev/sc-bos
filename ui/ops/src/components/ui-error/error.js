import {statusCodeToString} from '@/components/ui-error/util';
import {StatusCode} from 'grpc-web';
import {defineStore} from 'pinia';
import {computed, ref, watch} from 'vue';

/**
 * @typedef {Object} UiError
 * @property {string} name
 * @property {Error} source
 * @property {number} timestamp
 * @property {string} id
 */

export const useErrorStore = defineStore('error', () => {
  const _errorMap = ref({});
  let _id = 0;

  /**
   * @param {RemoteResource} resource
   * @param {RemoteError} error
   */
  function addError(resource, error) {
    const e = /** @type {UiError} */{name: error.name, source: error.error, timestamp: Date.now(), id: _id++};
    console.error(`[${(new Date(e.timestamp)).toLocaleTimeString()}] ${e.name} | ${statusCodeToString(e.source.code)}: ${e.source.message}`, e.source);

    _errorMap.value[e.id] = e;
    // auto-clear errors after 1 minute
    setTimeout(() => {
      clearError(e);
    }, 60 * 1000);
  }

  /**
   *
   * @param {UiError} error
   */
  function clearError(error) {
    delete(_errorMap.value[error.id]);
  }

  /**
   * @return {UiError[]}
   */
  const errors = computed(() => {
    return Object.values(_errorMap.value);
  });

  /**
   * @param {ActionTracker<?>} actionTracker
   * @return {WatchStopHandle}
   */
  function registerTracker(actionTracker) {
    if (Object.hasOwn(actionTracker, 'error')) {
      return watch(() => actionTracker.error, (error) => {
        if (error && error.code !== StatusCode.OK) {
          addError(actionTracker, error);
        }
      });
    } else if (actionTracker.value) {
      return watch(() => actionTracker.value.error, (error) => {
        if (error && error.code !== StatusCode.OK) {
          addError(actionTracker.value, error);
        }
      });
    }
    return () => {}; // shouldn't happen
  }

  /**
   *
   * @param {ResourceValue<?, ?>} resourceValue
   * @return {WatchStopHandle}
   */
  function registerValue(resourceValue) {
    return watch(() => resourceValue.streamError, (error) => {
      if (error && error.code !== StatusCode.OK) {
        addError(resourceValue, error);
      }
    });
  }

  /**
   *
   * @param {ResourceCollection<?, ?>} collection
   * @return {WatchStopHandle}
   */
  function registerCollection(collection) {
    if (Object.hasOwn(collection, 'streamError')) {
      return watch(() => collection.streamError, (error) => {
        if (error && error.code !== StatusCode.OK) {
          addError(collection, error);
        }
      });
    } else {
      return watch(() => collection.resources?.streamError, (error) => {
        if (error && error.code !== StatusCode.OK) {
          addError(collection.resources, error);
        }
      });
    }
  }

  return {
    addError,
    clearError,
    errors,
    registerTracker,
    registerValue,
    registerCollection
  };
});
