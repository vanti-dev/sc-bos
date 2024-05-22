import {defineStore} from 'pinia';
import {ref} from 'vue';

export const useSidebarStore = defineStore('sidebar', () => {
  // RIGHT SIDEBAR //
  const visible = ref(false);
  const data = ref({});
  const title = ref('');
  // for use when targeting a specific node
  const sidebarNode = ref({name: ''});

  const listedDevice = ref({});

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
    data.value = {};
    listedDevice.value = {};
    title.value = '';
    visible.value = false;
  };

  return {
    // RIGHT SIDEBAR
    visible,
    data,
    title,
    sidebarNode,
    listedDevice,
    toggleSidebar,
    closeSidebar,
    resetSidebarToDefaults
  };
});
