import {useCohortHealthStore, useCohortStore} from '@/stores/cohort.js';
import {acceptHMRUpdate, defineStore} from 'pinia';
import {ref} from 'vue';

// use config stores settings the user has chosen during their interaction with the ui.
export const useUserConfig = defineStore('userConfig', () => {
  // The SC node we're getting services from, if absent get services from the node we're communicating directly with.
  const node = ref(null);

  return {
    node
  };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useUserConfig, import.meta.hot));
}
