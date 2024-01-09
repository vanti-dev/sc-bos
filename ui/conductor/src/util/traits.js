import {watch} from 'vue';
import {closeResource} from '@/api/resource';
import {deepEqual} from 'vuetify/src/util/helpers';

/**
 * Watches specified props and performs resource-related actions based on changes.
 *
 * @template T
 * @param {(() => T) | (() => T)[]} watchedProps - function/array of fns representing the watched props.
 * @param {RemoteResource<T> | RemoteResource<T>[]} resourceValue - resource(s) to be closed and pulled.
 * @param {(name: object, resource: T) => void} apiCalls - The function to pull a single
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
export const watchResource = (watchedProps, resourceValue, apiCalls) => {
  watch(
      watchedProps,
      (newProps, oldProps) => {
        // If resourceValue is not an array, make it an array
        const resources = Array.isArray(resourceValue) ? resourceValue : [resourceValue];
        // Check if the props have changed
        const oldNewEqual = deepEqual(newProps, oldProps);

        // If the props haven't changed, do nothing
        if (oldNewEqual) return;

        // If the props have changed (either the name or the paused state or both), close the old resources
        // If the resources is an array, loop through and close all resources
        resources.forEach((resource) => closeResource(resource));

        if (!newProps?.paused) { // If not paused, pull new resource
          // Pull new resource
          resources.forEach((resource) => apiCalls({name: newProps?.name}, resource));
        }
      },
      {immediate: true, deep: true, flush: 'sync'}
  );
};
