import {closeResource} from '@/api/resource';
import {toValue} from '@/util/vue';
import {watch} from 'vue';
import {deepEqual} from 'vuetify/src/util/helpers';

/**
 * Calls apiCalls each time name changes, tracking and managing resource cleanup for you.
 *
 * @param {MaybeRefOrGetter<string>} name - string representing the name of the device
 * @param {MaybeRefOrGetter<boolean>} paused - boolean representing whether the data stream is paused
 * @param {(name: string) => ResourceValue<any, any>} apiCalls - array of functions that return a resource
 * @example
 * watchResource(
 *   () => props.name,
 *   () => props.paused,
 *   (name) => {
 *     pullAirTemperature({name}, resource);
 *     return resource;
 *   }
 * );
 */
export const watchResource = (name, paused, ...apiCalls) => {
  let resources = [];

  watch(
      [() => toValue(name), () => toValue(paused)],
      (newProps, oldProps) => {
        // Check if the props have changed
        const oldNewEqual = deepEqual(newProps, oldProps);

        // If the props haven't changed, do nothing
        if (oldNewEqual) return;
        // Separate the array values into readable variables
        const name = newProps[0];
        const paused = newProps[1];

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
