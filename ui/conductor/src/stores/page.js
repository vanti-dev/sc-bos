import {defineStore} from 'pinia';
import {computed, ref, watch} from 'vue';
import {useRoute} from 'vue-router/composables';
import {camelCasingString} from '@/util/string';

export const usePageStore = defineStore('page', () => {
  const currentRoute = useRoute();

  const pageRoute = ref('');
  const subPageRoute = ref('');
  const pageType = ref(
      /** @type {Object.<string, boolean>} */
      {
        automations: false,
        devices: false,
        editorMode: false,
        ops: false,
        site: false,
        system: false
      }
  );
  const subPageType = ref(
      /** @type {Object.<string, boolean>} */
      {}
  );
  const showSidebar = ref(false);
  const sidebarData = ref({});
  const sidebarTitle = ref('');
  // for use when targeting a specific node
  const sidebarNode = ref({name: ''});

  const requiredSlots = computed(() => {
    let slots;
    if (pageType.value.automations || pageType.value.system) {
      slots = ['actions'];
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
    showSidebar.value = false;
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
    const subPage = toPath.substring(secondToSlashIndex + 1);

    pageRoute.value = page;
    subPageRoute.value = subPage;
  }


  // Separating the path into main and sub values
  watch(
      () => currentRoute.fullPath,
      (newPath, oldPath) => {
        setActivePage(newPath);
      },
      {immediate: true, deep: true, flush: 'sync'}
  );

  // Setting Page route
  watch(
      pageRoute,
      (newRoute, oldRoute) => {
        if (camelCasingString(oldRoute)) pageType.value[camelCasingString(oldRoute)] = false;
        pageType.value[camelCasingString(newRoute)] = true;

        if (newRoute) closeSidebar();
      },
      {immediate: true, deep: true, flush: 'sync'}
  );

  // Setting Sub page route (sidebar nav values)
  watch(
      subPageRoute,
      (newRoute, oldRoute) => {
        if (camelCasingString(oldRoute)) subPageType.value[camelCasingString(oldRoute)] = false;
        subPageType.value[camelCasingString(newRoute)] = true;

        if (newRoute) closeSidebar();
      },
      {immediate: true, deep: true, flush: 'sync'}
  );

  return {
    pageRoute,
    subPageRoute,
    pageType,
    subPageType,
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
