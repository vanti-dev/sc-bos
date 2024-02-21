import {closeResource} from '@/api/resource';
import {toValue} from '@/util/vue';
import {computed, watch} from 'vue';
import {deepEqual} from 'vuetify/src/util/helpers';

/**
 * Converts a ref or reactive object to a computed property that evaluates to a query object.
 *
 * @template {{name: string}} T
 * @param {MaybeRefOrGetter<string|T|null>} input - input value to be converted into a query object.
 * @return {
 *  import('vue').ComputedRef<T|null>
 * } - computed property that evaluates to a query object.
 */
export const toQueryObject = (input) => {
  return computed(() => {
    const inputValue = toValue(input);

    // If no input present, return null
    if (!inputValue) return null;

    // If input is a string, return an object with the name property
    if (typeof inputValue === 'string') {
      return {name: inputValue};
      //
      // If input is an object, return as is
    } else {
      return inputValue;
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
        const req = newSource[0];
        const paused = newSource[1];

        // If the props have changed (either the name or the paused state or both), close the old resources
        // and empty the array
        resources.forEach((resource) => closeResource(resource));
        resources = [];
        if (!paused) { // If not paused, pull new resource
          apiCalls.forEach((apiCall) => resources.push(apiCall(req)));
        }
      },
      {immediate: true, deep: true, flush: 'sync'}
  );
};
