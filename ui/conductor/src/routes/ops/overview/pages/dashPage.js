import {builtinLayouts} from '@/layout/builtinLayouts.js';
import {useUiConfigStore} from '@/stores/ui-config.js';
import {findActiveItem} from '@/util/router.js';
import {toValue} from '@/util/vue.js';
import {builtinWidgets} from '@/widgets/pallet.js';
import {computed} from 'vue';

/**
 * @param {MaybeRefOrGetter<string|string[]>} path
 * @return {{
 *   layout: import('vue').Component,
 *   main: ComputedRef<{component:*,props:Object}[]>,
 *   after: ComputedRef<{component:*,props:Object}[]>,
 *   title: ComputedRef<string>
 * }}
 */
export default function useDashPage(path) {
  const uiConfig = useUiConfigStore();

  // uri encoded path segments from path.
  const pathSegments = computed(() => {
    let p = toValue(path);
    if (!Array.isArray(p)) p = p.split('/');
    return p.map(s => encodeURIComponent(s));
  });

  const pageConfigRoot = computed(
      () => uiConfig.getOrDefault('ops.overview', {}));
  // the current page config, detailing the title, icon and widgets to show.
  const pageConfig = computed(() => {
    if (!pathSegments.value.length || !pageConfigRoot.value?.children?.length) return null;
    // Note: findActiveItem takes uri-encoded path segments
    return findActiveItem(pageConfigRoot.value.children, pathSegments.value);
  });

  // config for all the widgets we support
  const powerHistoryConfig = computed(() => {
    const v = pageConfig.value?.widgets?.showEnergyConsumption;
    if (!v) return false;
    return {
      demand: /** @type {string} */ v,
      hideTotal: true,
      chartTitle: 'Energy Consumption'
    };
  });
  const zoneNotificationConfig = computed(() => {
    const v = pageConfig.value?.widgets?.showNotifications;
    if (!v) return false;
    return {
      zone: /** @type {string} */ v
    };
  });
  const presenceConfig = computed(() => {
    const v = pageConfig.value?.widgets?.showOccupancy;
    if (!v) return false;
    return {
      source: /** @type {string} */ v
    };
  });
  const environmentConfig = computed(() => {
    const v = pageConfig.value?.widgets?.showEnvironment;
    if (!v) return false;
    return {
      internal: /** @type {string} */ v.indoor,
      external: /** @type {string} */ v.outdoor
    };
  });

  const title = computed(() => pageConfig.value?.title ?? '');
  const extendedTitle = computed(() => {
    if (!title.value) return '';
    return title.value + ' Status Overview';
  });

  const addIfPresent = (arr, props, component) => {
    props = toValue(props);
    if (props) {
      arr.push({props, component});
    }
  };
  const mainWidgets = computed(() => {
    const res = [];
    addIfPresent(res, powerHistoryConfig, builtinWidgets['power-history']);
    addIfPresent(res, zoneNotificationConfig, builtinWidgets['zone-notifications']);
    return res;
  });
  const afterWidgets = computed(() => {
    const res = [];
    addIfPresent(res, presenceConfig, builtinWidgets['presence']);
    addIfPresent(res, environmentConfig, builtinWidgets['environmental']);
    return res;
  });

  return {
    layout: builtinLayouts['main-side'],
    title: extendedTitle,
    main: mainWidgets,
    after: afterWidgets
  };
}
