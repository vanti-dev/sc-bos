import {defineStore} from 'pinia';
import {reactive} from 'vue';

export const useIntersectedItemsStore = defineStore('intersectedItems', () => {
  const intersectedItemNames = reactive(
      /** @type {Object.<string, number>} */
      {}
  );

  /**
   *
   * @param {IntersectionObserverEntry[]} entries
   * @param {IntersectionObserver} observer
   * @param {string} name
   */
  const intersectionHandler = (entries, observer, name) => {
    const anyIntersect = entries.some((entry) => entry.isIntersecting);
    if (anyIntersect) {
      createName(name);
      //
    } else {
      clearName(name);
    }
  };

  const hasName = (name) => {
    return intersectedItemNames[name] > 0;
  };

  const clearName = (name) => {
    if (hasName(name)) {
      intersectedItemNames[name]--;
      if (!hasName(name)) {
        delete(intersectedItemNames[name]);
      }
    }
  };

  const createName = (name) => {
    if (hasName(name)) {
      intersectedItemNames[name]++;
    } else {
      intersectedItemNames[name] = 1;
    }
  };

  return {
    intersectedItemNames,

    intersectionHandler,
    hasName,
    clearName,
    createName
  };
});
