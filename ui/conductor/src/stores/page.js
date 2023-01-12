import {defineStore} from 'pinia';
import {ref} from 'vue';

export const usePageStore = defineStore('page', () => {
  const showSidebar = ref(false);
  const sidebarData = ref({});
  const sidebarTitle = ref('');

  /**
   *
   */
  function toggleSidebar() {
    showSidebar.value = !showSidebar.value;
  }

  return {
    showSidebar,
    sidebarData,
    sidebarTitle,
    toggleSidebar
  };
});
