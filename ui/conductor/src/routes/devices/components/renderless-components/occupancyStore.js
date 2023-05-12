import {defineStore} from 'pinia';
import {del, reactive, set} from 'vue';
import {closeResource} from '@/api/resource';
import {pullOccupancy} from '@/api/sc/traits/occupancy';

export const useOccupancyStore = defineStore('occupancy', () =>{
  //
  //
  // Methods
  /**
   *
   * @param {Device.AsObject} item
   * @return {undefined|occupancyTrait}
   */
  function findOccupancySensor(item) {
    const occupancyTrait = item.metadata.traitsList.find(trait => {
      if (trait.name.includes('Occupancy')) return trait;
    });

    if (occupancyTrait) return occupancyTrait;
    else return undefined;
  }


  // ///////////////////
  //
  // Intersection
  const intersectedItemNames = reactive(/** @type{[{[string]: boolean}]} */{});

  /**
   *
   * @param {IntersectionObserverEntry} entries
   * @param {IntersectionObserver} observer
   * @param {string} name
   */
  const onRowIntersect = (entries, observer, name) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        // when changing table page and name exist
        set(intersectedItemNames, name, true);
        //
      } else del(intersectedItemNames, name);
    });
  };

  //
  // Action
  // Let's see if the row is intersecting
  const handleStream = (name, paused, occupancyValue) => {
    closeResource(occupancyValue);

    if (!name || paused) {
      return;
    }

    pullOccupancy(name, occupancyValue);
  };

  const resetIntersectedItemNames = () => {
    Object.keys(intersectedItemNames).forEach(key => del(intersectedItemNames, key));
  };

  return {
    findOccupancySensor,

    // Intersection
    intersectedItemNames,
    onRowIntersect,
    handleStream,
    resetIntersectedItemNames
  };
});
