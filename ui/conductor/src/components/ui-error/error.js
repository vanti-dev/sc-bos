import {statusCodeToString} from '@/components/ui-error/util';
import {StatusCode} from 'grpc-web';
import {defineStore} from 'pinia';
import Vue, {computed, ref, watch} from 'vue';

import {useAccountStore} from '@/stores/account';

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
    // eslint-disable-next-line max-len
    console.error(`[${(new Date(e.timestamp)).toLocaleTimeString()}] ${e.name} | ${statusCodeToString(e.source.code)}: ${e.source.message}`, e.source);

    const account = useAccountStore();
    // Log the user out if we get a permission denied error
    // and clear the local storage
    if (statusCodeToString(e.source.code) === 'PERMISSION DENIED') {
      account.logout();
      localStorage.clear();
    }

    Vue.set(_errorMap.value, e.id, e);
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
    Vue.delete(_errorMap.value, error.id);
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
    if (actionTracker.hasOwnProperty('error')) {
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
    if (collection.hasOwnProperty('streamError')) {
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
