import {defineStore} from 'pinia';
import {ref} from 'vue';

export const useServicesStore = defineStore('services', () => {
  // The SC node we're getting services from, if absent get services from the node we're communicating directly with.
  const node = ref(null);

  return {
    node
  };
});
