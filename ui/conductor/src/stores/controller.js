import {defineStore} from 'pinia';
import {ref} from 'vue';

export const useControllerStore = defineStore('controller', () => {
  // todo: get this from somewhere
  const controllerName = ref(''); // blank means use the controllers default name
  return {
    controllerName
  };
});
