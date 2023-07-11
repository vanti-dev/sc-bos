import {unref} from 'vue';

/**
 * @typedef MaybeRefOrGetter
 * @template T
 * @type {T | import('vue').Ref<T> | (() => T)}
 */

/**
 * Return the concrete value from MaybeRefOrGetter.
 *
 * @param {MaybeRefOrGetter<T>} source
 * @return {T}
 * @template T
 */
export function toValue(source) {
  return typeof source === 'function' ? source() : unref(source);
}
