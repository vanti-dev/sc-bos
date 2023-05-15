import {defineStore} from 'pinia';
import {del, reactive, ref, set, watch} from 'vue';
import {closeResource} from '@/api/resource';
import {pullOccupancy} from '@/api/sc/traits/occupancy';

export const useOccupancyStore = defineStore('occupancy', () =>{
  //
  // State
  const activePage = ref(1);
  const devicesPerPage = ref(10);
  const intersectedItemNames = reactive(/** @type{[{[string]: boolean}]} */{});
  const perPageChoices = [
    {text: '5', value: 5},
    {text: '10', value: 10},
    {text: '20', value: 20},
    {text: 'All', value: -1}
  ];

  //
  // Actions
  /**
   *
   * @param {Device.AsObject} item
   * @param {string} type
   * @return {undefined|trait}
   */
  const findSensor = (item, type) => {
    const trait = item.metadata.traitsList.find(trait => {
      if (trait.name.includes(type)) return trait;
    });

    if (trait) return trait;
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


  const handleStream = (name, paused, value) => {
    closeResource(value);

    if (!name || paused) {
      return;
    }

    pullOccupancy(name, value);
  };

  // stopping left over streams
  const resetIntersectedItemNames = () => {
    Object.keys(intersectedItemNames).forEach((key, index) => {
      del(intersectedItemNames, key);
    });
  };


  //
  //
  // Watchers

  let timeoutId;
  // 3 seconds timeout to reset go to page input to active page
  watch(activePage, (newPage, oldPage) => {
    clearTimeout(timeoutId);

    // if we specified the new page
    if (newPage && oldPage !== '' && newPage !== oldPage) resetIntersectedItemNames();

    // if we have an empty input
    if (newPage === '' || isNaN(newPage)) {
      timeoutId = setTimeout(() => {
        activePage.value = Number(oldPage);
      }, 3000);
    }
  });

  // stopping left over streams on device per page value change
  watch(devicesPerPage, value => {
    if (value !== -1 && Object.keys(intersectedItemNames).length > value) {
      resetIntersectedItemNames();
    }
  });

  return {
    activePage,
    devicesPerPage,
    perPageChoices,
    findSensor,

    // Intersection
    intersectedItemNames,
    onRowIntersect,
    handleStream,
    resetIntersectedItemNames
  };
});
