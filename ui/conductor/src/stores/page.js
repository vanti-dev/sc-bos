import {defineStore} from 'pinia';
import {ref} from 'vue';

export const usePageStore = defineStore('page', () => {
  const showSidebar = ref(false);
  const sidebarData = ref({});
  const sidebarTitle = ref('');

  // for use when targeting a specific node
  const sidebarNode = ref({name: ''});

  /**
   *
   */
  function toggleSidebar() {
    showSidebar.value = !showSidebar.value;
  }

  /**
   *
   */
  function closeSidebar() {
    toggleSidebar();
    sidebarTitle.value = '';
    sidebarData.value = {};
  };

  return {
    showSidebar,
    sidebarData,
    sidebarTitle,
    toggleSidebar,
    closeSidebar,
    sidebarNode
  };
});
