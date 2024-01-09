import {toValue} from '@/util/vue';
import {watch} from 'vue';
import {closeResource} from '@/api/resource';
import {deepEqual} from 'vuetify/src/util/helpers';

/**
 * Watches specified props and performs resource-related actions based on changes.
 *
 * @template T
 * @param {MaybeGetterOrRef<string>} name - function/array of fns representing the watched props.
 * @param {MaybeGetterOrRef<boolean>} paused - function/array of fns representing the watched props.
 * @param {Array<(name: string) => ResourceValue>} apiCalls - The function to pull a single
 * or multiple resource data.
 * @example
 * watchResource(
 *   [() => props.paused, () => props.name],
 *   airTemperatureResource,
 *   (name, resource) => {
 *     pullAirTemperature(name, resource);
 *   }
 * );
 */
export const watchResource = (name, paused, apiCalls) => {
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
        // If the resources is an array, loop through and close all resources
        resources.forEach((resource) => closeResource(resource));
        resources = [];

        if (!paused) { // If not paused, pull new resource
          apiCalls.forEach((apiCall) => resources.push(apiCall(name)));
        }
      },
      {immediate: true, deep: true, flush: 'sync'}
  );
};
