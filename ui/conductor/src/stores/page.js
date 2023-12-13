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
   *
   */
  function toggleSidebar() {
    showSidebar.value = !showSidebar.value;
  }

  /**
   *
   */
  function closeSidebar() {
    showSidebar.value = false;
    sidebarData.value = {};
    listedDevice.value = {};
    sidebarTitle.value = '';
  };

  //
  //
  // LEFT NAVIGATION SIDEBAR
  const drawer = ref(true);
  const miniVariant = ref(true);
  const drawerWidth = ref(45);
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

    // LEFT NAVIGATION SIDEBAR
    drawer,
    miniVariant,
    drawerWidth,
    pinDrawer
  };
});
