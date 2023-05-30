import {defineStore} from 'pinia';
import {reactive} from 'vue';

export const useIntersectedItemsStore = defineStore('intersectedItems', () => {
  const intersectedItemNames = reactive(
      /** @type {Array<Object.<string, boolean>>} */
      []
  );

  /**
   *
   * @param {IntersectionObserverEntry} entries
   * @param {IntersectionObserver} observer
   * @param {string} table
   * @param {string} name
   */
  const intersectionHandler = (entries, observer, table, name) => {
    entries.forEach((entry, index) => {
      if (entry.isIntersecting) {
        set(intersectedItemNames[table], name, true);
        //
      } else {
        del(intersectedItemNames[table], name);
      }
    });
  };

  return {
    intersectedItemNames,

    intersectionHandler
  };
});
