import {newActionTracker} from '@/api/resource';
import {getMetadata} from '@/api/sc/traits/metadata';
import {defineStore} from 'pinia';
import {ref} from 'vue';

export const useControllerStore = defineStore('controller', () => {
  const controllerName = ref(''); // blank means use the controllers default name

  /**
   * @return {Promise<void>}
   */
  async function sync() {
    // we only need to do this once
    if (controllerName.value === '') {
      const meta = await getMetadata({}, newActionTracker());
      controllerName.value = meta.name;
    }
  }

  return {
    controllerName,
    sync
  };
});
