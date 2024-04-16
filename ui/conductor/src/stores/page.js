import {defineStore} from 'pinia';
import {ref} from 'vue';

export const usePageStore = defineStore('page', () => {
  // RIGHT SIDEBAR //
  const showSidebar = ref(false);
  const sidebarData = ref({});
  const sidebarTitle = ref('');
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
    sidebarTitle.value = '';
    showSidebar.value = false;
  };

  //
  //
  // LEFT NAVIGATION SIDEBAR
  const drawer = ref(true);
  const miniVariant = ref(true);
  const drawerWidth = ref(60);
  const pinDrawer = ref(false);

  return {
    // RIGHT SIDEBAR
    showSidebar,
    sidebarData,
    sidebarTitle,
    sidebarNode,
    listedDevice,
    toggleSidebar,
    closeSidebar,
    resetSidebarToDefaults,

    // LEFT NAVIGATION SIDEBAR
    drawer,
    miniVariant,
    drawerWidth,
    pinDrawer
  };
});
