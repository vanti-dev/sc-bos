import {defineStore} from 'pinia';
import {del, reactive, set} from 'vue';
import {closeResource} from '@/api/resource';
import {pullOccupancy} from '@/api/sc/traits/occupancy';

export const useOccupancyStore = defineStore('occupancy', () =>{
  //
  // State
  const intersectedItemNames = reactive(/** @type{[{[string]: boolean}]} */{});


  //
  // Actions
  /**
   *
   * @param {Device.AsObject} item
   * @return {undefined|occupancyTrait}
   */
  const findOccupancySensor = (item) => {
    const occupancyTrait = item.metadata.traitsList.find(trait => {
      if (trait.name.includes('Occupancy')) return trait;
    });

    if (occupancyTrait) return occupancyTrait;
    else return undefined;
  };


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
