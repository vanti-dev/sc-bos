import {closeResource} from '@/api/resource';
import {toValue} from '@/util/vue';
import {computed, watch} from 'vue';
import {deepEqual} from 'vuetify/src/util/helpers';

/**
 * Converts a ref or reactive object to a computed property that evaluates to a query object.
 *
 * @template T
 * @param {MaybeRefOrGetter<T>} input - input value to be converted into a query object.
 * @return {
 *  import('vue').ComputedRef<{name: string} | T | null>
 * } - computed property that evaluates to a query object.
 */
export const toQueryObject = (input) => {
  return computed(() => {
    // If no input present, return null
    if (!toValue(input)) return null;

    // If input is a string, return an object with the name property
    if (typeof toValue(input) === 'string') {
      return {name: toValue(input)};
      //
      // If input is an object, return the object
    } else {
      return toValue(input);
    }
  });
};

/**
 * Calls apiCalls each time the query changes, tracking and managing resource cleanup.
 *
 * @template T
 * @param {MaybeRefOrGetter<T>} query - object representing the request to the API
 * @param {MaybeRefOrGetter<boolean>} [paused] - boolean representing whether the data stream is paused
 * @param {(req: T) => T<any, any>} apiCalls - array of functions that return a resource
 * @example
 * watchResource(
 *   () => toValue(toQueryObject(query)),
 *   () => toValue(paused),
 *   (req) => {
 *     pullAirTemperature(req, resource);
 *     return resource;
 *   }
 * );
 */
export const watchResource = (query, paused = false, ...apiCalls) => {
  let resources = [];

  watch(
      [() => toValue(query), () => toValue(paused)],
      (newSource, oldSource) => {
        // Check if the props have changed
        const oldNewEqual = deepEqual(newSource, oldSource);

        // If the props haven't changed, do nothing
        if (oldNewEqual) return;
        // Separate the array values into readable variables
        const name = newSource[0];
        const paused = newSource[1];

        // If the props have changed (either the name or the paused state or both), close the old resources
        // and empty the array
        resources.forEach((resource) => closeResource(resource));
        resources = [];
        if (!paused) { // If not paused, pull new resource
          apiCalls.forEach((apiCall) => resources.push(apiCall(name)));
        }
      },
      {immediate: true, deep: true, flush: 'sync'}
  );
};
