import {defineStore} from 'pinia';
import {ref} from 'vue';

export const useSidebarStore = defineStore('sidebar', () => {
  // RIGHT SIDEBAR //
  const showSidebar = ref(false);
  const sidebarData = ref({});
  const title = ref('');
  // for use when targeting a specific node
  const sidebarNode = ref({name: ''});

  const listedDevice = ref({});

  /**
   * Open or close sidebar
   */
  const toggleSidebar = () => {
    showSidebar.value = !showSidebar.value;
  };

  /**
   * Close sidebar if visible and reset sidebar data
   */
  const closeSidebar = () => {
    showSidebar.value = false;
    resetSidebarToDefaults();
  };

  /**
   * Reset the sidebar data to default values
   */
  const resetSidebarToDefaults = () => {
    sidebarData.value = {};
    listedDevice.value = {};
    title.value = '';
    showSidebar.value = false;
  };

  return {
    // RIGHT SIDEBAR
    showSidebar,
    sidebarData,
    title,
    sidebarNode,
    listedDevice,
    toggleSidebar,
    closeSidebar,
    resetSidebarToDefaults
  };
});
