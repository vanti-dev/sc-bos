import {defineStore} from 'pinia';
import {ref} from 'vue';

export const useControllerStore = defineStore('controller', () => {
  // todo: get this from somewhere
  const controllerName = ref('test-ac');
  return {
    controllerName
  };
});
