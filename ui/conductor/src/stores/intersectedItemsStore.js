import {defineStore} from 'pinia';
import {del, reactive, set} from 'vue';

export const useIntersectedItemsStore = defineStore('intersectedItems', () => {
  const intersectedItemNames = reactive(
      /** @type {Object.<string, boolean>} */
      {}
  );

  /**
   *
   * @param {IntersectionObserverEntry} entries
   * @param {IntersectionObserver} observer
   * @param {string} name
   */
  const intersectionHandler = (entries, observer, name) => {
    entries.forEach((entry, index) => {
      if (entry.isIntersecting) {
        createName(name);
        //
      } else {
        clearName(name);
      }
    });
  };

  const clearName = (name) => {
    del(intersectedItemNames, name);
  };

  const createName = (name) => {
    const match = Object.keys(intersectedItemNames).find((itemName) => itemName) === name;

    if (match) {
      intersectedItemNames[name]++;
    } else {
      set(intersectedItemNames, name, 1);
    }
  };

  return {
    intersectedItemNames,

    intersectionHandler,
    clearName,
    createName
  };
});
