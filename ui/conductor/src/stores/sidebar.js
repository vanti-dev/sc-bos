import {defineStore} from 'pinia';
import {ref} from 'vue';

export const useSidebarStore = defineStore('sidebar', () => {
  // indicates whether the sidebar is visible or not
  const visible = ref(false);
  // The title of the sidebar
  const title = ref('');
  // component provided data used to communicate between the main page and the sidebar
  const data = ref({});

  /**
   * Open or close sidebar
   */
  const toggleSidebar = () => {
    visible.value = !visible.value;
  };

  /**
   * Close sidebar if visible and reset sidebar data
   */
  const closeSidebar = () => {
    visible.value = false;
    resetSidebarToDefaults();
  };

  /**
   * Reset the sidebar data to default values
   */
  const resetSidebarToDefaults = () => {
    title.value = '';
    data.value = {};
    visible.value = false;
  };

  return {
    visible,
    title,
    data,
    toggleSidebar,
    closeSidebar,
    resetSidebarToDefaults
  };
});
