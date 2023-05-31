import {Device} from '@sc-bos/ui-gen/proto/devices_pb';

/**
 *
 * @param {Device.AsObject} item
 * @param {string} trait
 * @return {Array<string>}
 */
export function hasTrait(item, trait) {
  const traitsArray = item?.metadata?.traitsList.map(trait => {
    return trait.name;
  });

  const sensors = traitsArray.map((item) => {
    const arr = item.split('.');
    return arr.slice(2).join('.');
  });

  return sensors.includes(trait);
}
