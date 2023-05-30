import {closeResource} from '@/api/resource';
import {pullOccupancy} from '@/api/sc/traits/occupancy';

/**
 *
 * @param {string} name
 * @param {boolean} paused
 * @param {*} value
 */
export function handleOccupancyStream(name, paused, value) {
  if (!name || paused) {
    closeResource(value);
    return;
  }

  pullOccupancy(name, value);
};


/**
 *
 * @param {*} value
 */
export function closeStream(value) {
  closeResource(value);
}
