import {defineStore} from 'pinia';
import {ref} from 'vue';

export const usePageStore = defineStore('page', () => {
  const showSidebar = ref(false);
  const sidebarData = ref({});

  /**
   *
   */
  function toggleSidebar() {
    showSidebar.value = !showSidebar.value;
  }

  return {
    showSidebar,
    sidebarData,
    toggleSidebar
  };
});
