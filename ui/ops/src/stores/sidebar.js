import {defineStore} from 'pinia';
import {ref, shallowRef, watch} from 'vue';

export const useSidebarStore = defineStore('sidebar', () => {
  // indicates whether the sidebar is visible or not
  const visible = ref(false);
  // The title of the sidebar
  const title = ref('');
  // component provided data used to communicate between the main page and the sidebar
  const data = ref({});
  // a dynamic component used for the sidebar
  const component = shallowRef(null);

  // todo: remove this once we've created a better way to show sidebars.
  watch(data, () => {
    component.value = null;
  }, {flush: 'sync'});

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
    resetSidebarToDefaults();
  };

  /**
   * Reset the sidebar data to default values
   */
  const resetSidebarToDefaults = () => {
    visible.value = false;
    title.value = '';
    data.value = {};
    component.value = null;
  };

  return {
    visible,
    title,
    data,
    component,
    toggleSidebar,
    closeSidebar,
    resetSidebarToDefaults
  };
});
