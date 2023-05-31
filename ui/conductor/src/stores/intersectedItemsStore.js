import {defineStore} from 'pinia';
import {del, reactive, set} from 'vue';

export const useIntersectedItemsStore = defineStore('intersectedItems', () => {
  const intersectedItemNames = reactive(
      /** @type {Object.<string, number>} */
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

  const hasName = (name) => {
    return intersectedItemNames[name] > 0;
  };

  const clearName = (name) => {
    if (hasName(name)) {
      intersectedItemNames[name]--;
      if (!hasName(name)) {
        del(intersectedItemNames, name);
      }
    }
  };

  const createName = (name) => {
    if (hasName(name)) {
      intersectedItemNames[name]++;
    } else {
      set(intersectedItemNames, name, 1);
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
