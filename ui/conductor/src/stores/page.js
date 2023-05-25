import {defineStore} from 'pinia';
import {computed, ref, watch} from 'vue';
import {useRoute} from 'vue-router/composables';

export const usePageStore = defineStore('page', () => {
  const currentRoute = useRoute();
  const pageRoute = ref('');
  const pageType = ref(
      /** @type {Object.<string, boolean>} */
      {
        automations: false,
        devices: false,
        editorMode: false,
        site: false,
        system: false
      }
  );
  const showSidebar = ref(false);
  const sidebarData = ref({});
  const sidebarTitle = ref('');
  // for use when targeting a specific node
  const sidebarNode = ref({name: ''});

  const requiredSlots = computed(() => {
    let slots;
    if (pageType.value.automations || pageType.value.system) {
      slots = ['active', 'actions'];
    } else if (pageType.value.devices) {
      slots = ['hotpoints'];
    } else slots = [];

    return slots;
  });

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
  }

  /**
   *
   * @param {string} toPath
   */
  function setActivePage(toPath) {
    const firstToSlashIndex = toPath.indexOf('/');
    const secondToSlashIndex = toPath.indexOf('/', firstToSlashIndex + 1);

    const page = toPath.substring(firstToSlashIndex + 1, secondToSlashIndex);
    pageRoute.value = page;
  }

  watch(
      () => currentRoute.fullPath,
      (newPath, oldPath) => {
        setActivePage(newPath);
      },
      {immediate: true, deep: true, flush: 'sync'}
  );

  watch(
      pageRoute,
      (newRoute, oldRoute) => {
        if (oldRoute) pageType.value[oldRoute] = false;
        pageType.value[newRoute] = true;
      },
      {immediate: true, deep: true, flush: 'sync'}
  );

  return {
    pageRoute,
    pageType,
    showSidebar,
    sidebarData,
    sidebarTitle,
    sidebarNode,

    requiredSlots,

    toggleSidebar,
    closeSidebar,
    setActivePage
  };
});
