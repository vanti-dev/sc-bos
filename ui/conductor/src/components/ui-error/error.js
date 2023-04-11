import {defineStore} from 'pinia';
import Vue, {computed, ref} from 'vue';

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
    Vue.set(_errorMap.value, e.id, e);
    // todo: expire errors
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

  return {
    addError,
    clearError,
    errors
  };
});
