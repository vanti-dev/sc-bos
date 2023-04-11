import {defineStore} from 'pinia';
import Vue, {computed, ref, watch} from 'vue';
import {StatusCode} from 'grpc-web';
import {statusCodeToString} from '@/components/ui-error/util';

/**
 * @typedef {Object} UiError
 * @property {RemoteResource} resource
 * @property {Error} source
 * @property {number} timestamp
 * @property {string} id
 */

export const useErrorStore = defineStore('error', () => {
  const _errorMap = ref({});
  let _id = 0;

  /**
   * @param {RemoteResource} resource
   * @param {Error} error
   */
  function addError(resource, error) {
    const e = /** @type {UiError} */{resource, source: error, timestamp: Date.now(), id: _id++};
    // eslint-disable-next-line max-len
    console.error(`[${(new Date(e.timestamp)).toLocaleTimeString()}] ${statusCodeToString(e.source.code)}: ${e.source.message}`, e.source);
    Vue.set(_errorMap.value, e.id, e);
    // todo: auto-clear errors
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
   * @param {Ref<ActionTracker>} actionTracker
   * @return {WatchStopHandle}
   */
  function registerTracker(actionTracker) {
    return watch(() => actionTracker.value.error, (error) => {
      if (error && error.code !== StatusCode.OK) {
        addError(actionTracker.resource, error);
      }
    });
  }

  /**
   *
   * @param {Ref<ResourceValue>} resourceValue
   * @return {WatchStopHandle}
   */
  function registerValue(resourceValue) {
    return watch(() => resourceValue.value.streamError, (error) => {
      if (error && error.code !== StatusCode.OK) {
        addError(resourceValue.value, error);
      }
    });
  }

  /**
   *
   * @param {Ref<Collection>} collection
   * @return {WatchStopHandle}
   */
  function registerCollection(collection) {
    return watch(() => collection.value.resources.streamError, (error) => {
      if (error && error.code !== StatusCode.OK) {
        addError(collection.value.resources, error);
      }
    });
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
