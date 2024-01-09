import {watch} from 'vue';
import {closeResource} from '@/api/resource';
import {deepEqual} from 'vuetify/src/util/helpers';

/**
 * Watches specified props and performs resource-related actions based on changes.
 *
 * @param {Function[]} watchedProps - array of functions representing the watched props.
 * @param {any} resource - The resource to be closed and pulled.
 * @param {(params: object, resource: any) => void} apiCalls - The function to pull a single
 * or multiple resource data.
 * //
 * @example
 * watchResource(
 *   [() => props.paused, () => props.name],
 *   meterReadings,
 *   (params) => {
 *     pullMeterReading(params, meterReadings);
 *     describeMeterReading(params, meterReadingInfo);
 *   }
 * );
 */
export const watchResource = (watchedProps, resource, apiCalls) => {
  watch(
      watchedProps,
      ([newPaused, newName], [oldPaused, oldName]) => {
        // Check if request is the same
        const requestEqual = deepEqual(newName, oldName);
        // If name is a string, wrap it in an object
        const nameParam = typeof newName === 'string' ? {name: newName} : newName;

        // If it's paused and the request hasn't changed, do nothing
        if (newPaused === oldPaused && requestEqual) return;

        closeResource(resource); // Close existing resource first

        if (!newPaused) { // If not paused, pull new resource
          apiCalls(nameParam, resource); // Pull resource
        }
      },
      {immediate: true, deep: true, flush: 'sync'}
  );
};
