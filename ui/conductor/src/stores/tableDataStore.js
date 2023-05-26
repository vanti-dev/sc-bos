import {defineStore} from 'pinia';
import {del, reactive, ref, set, watch} from 'vue';
import {closeResource} from '@/api/resource';
import {pullOccupancy} from '@/api/sc/traits/occupancy';

export const useTableDataStore = defineStore('tableData', () => {
  //
  // State
  const activePage = ref(1);
  const itemsPerPage = ref(10);
  const triggerRerender = ref(0);
  const intersectedItemNames = reactive(
      /** @type {Object.<string, boolean>} */
      {}
  );
  const perPageChoices = [
    {text: '5', value: 5},
    {text: '10', value: 10},
    {text: '20', value: 20},
    {text: 'All', value: -1}
  ];

  const search = ref('');
  const tableSelection = ref([]);
  const withTopBar = ref({
    /** @type {Object.<string, boolean>} */
    automations: false,
    devices: true,
    ops: false,
    site: true,
    system: false
  });

  //
  // Actions
  /**
   *
   * @param {Device.AsObject} item
   * @param {string} type
   * @return {undefined|trait}
   */
  const findSensor = (item, type) => {
    const includedTrait = item.metadata.traitsList.find((trait) => {
      if (trait.name.includes(type)) return trait;
    });

    if (includedTrait) return includedTrait;
    else return undefined;
  };

  /**
   *
   * @param {IntersectionObserverEntry} entries
   * @param {IntersectionObserver} observer
   * @param {string} name
   */
  const intersectionHandler = (entries, observer, name) => {
    entries.forEach((entry, index) => {
      // console.log(Object.keys(intersectedItemNames) === entry.target.firstElementChild.innerText);

      if (entry.isIntersecting) {
        set(intersectedItemNames, name, true);
        //
      } else {
        del(intersectedItemNames, name);
      }
    });
  };

  /**
   *
   * @param {string} name
   * @param {boolean} paused
   * @param {*} value
   */
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
  // 3 seconds timeout to reset go to page input back to active page
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
  watch(itemsPerPage, () => {
    resetIntersectedItemNames();
    triggerRerender.value += 1;
  });

  return {
    activePage,
    itemsPerPage,
    triggerRerender,
    perPageChoices,
    findSensor,
    search,
    tableSelection,
    withTopBar,

    // Intersection
    intersectedItemNames,
    intersectionHandler,
    handleStream,
    resetIntersectedItemNames
  };
});
