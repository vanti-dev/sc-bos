import {builtinLayouts} from '@/dynamic/layout/pallet.js';
import useBuildingConfig from '@/routes/ops/overview/pages/buildingConfig.js';
import useDashPage from '@/routes/ops/overview/pages/dashPage.js';
import {useUiConfigStore} from '@/stores/uiConfig.js';
import {findActiveItem} from '@/util/router.js';
import {builtinWidgets} from '@/dynamic/widgets/pallet.js';
import {computed, markRaw, reactive, toValue} from 'vue';

/**
 * @param {MaybeRefOrGetter<string|string[]>} path - uri decoded path to the page
 * @return {{
 *   layout: ComputedRef<import('vue').Component>,
 *   config: ComputedRef<Object>,
 *   isLegacySubPage: ComputedRef<boolean>,
 *   isLegacyOverview: ComputedRef<boolean>,
 *   pageConfig: ComputedRef<undefined|*|OverviewChild[]|null>,
 *   pageConfigNorm: ComputedRef<{}>
 * }}
 */
export default function usePageConfig(path) {
  const uiConfig = useUiConfigStore();

  // uri encoded path segments from path.
  const pathSegments = computed(() => {
    let p = toValue(path);
    if (!Array.isArray(p)) p = p.split('/');
    return p.map(s => encodeURIComponent(s));
  });

  // The config object (in ops.paths) described by path.
  // If ops.pages is not defined - i.e. legacy config is being used - this will also be undefined.
  const pageConfig = computed(() => {
    const pages = uiConfig.getOrDefault('ops.pages');
    if (!pages) {
      return undefined; // fall back to legacy config
    }
    if (!pathSegments.value.length) {
      return pages[0];
    }
    return findActiveItem(pages, pathSegments.value);
  });
  // The same as pageConfig but with special values hydrated to their component references.
  const pageConfigNorm = computed(() => {
    const cfg = pageConfig.value;
    if (!cfg) return cfg;

    const normObj = (o) => {
      return Object.entries(o).reduce((acc, [k, v]) => {
        if (Array.isArray(v)) {
          acc[k] = v.map(o => normObj(o));
        } else if (typeof v === 'object') {
          acc[k] = normObj(v);
        } else if (typeof v === 'string') {
          if (v.startsWith('builtin:')) {
            const [, builtin] = v.split(':');
            switch (k) {
              case 'layout':
                acc[k] = markRaw(builtinLayouts[builtin]);
                return acc;
              case 'component':
                acc[k] = markRaw(builtinWidgets[builtin]);
                return acc;
            }
          }
          acc[k] = v;
        } else {
          acc[k] = v;
        }
        return acc;
      }, {});
    };

    return normObj(cfg);
  });

  const isLegacyOverview = computed(() => !pageConfig.value && pathSegments.value.length === 1);
  const isLegacySubPage = computed(() => !pageConfig.value && pathSegments.value.length > 1);

  // Page config taken from either the newer ops.pages or legacy ops.overview as needed.
  const configObj = computed(() => {
    if (isLegacyOverview.value) {
      return reactive(useBuildingConfig());
    }
    if (isLegacySubPage.value) {
      return reactive(useDashPage(() => pathSegments.value
          // skip /building prefix in segments as useDashPage isn't expecting it
          .slice(1)
          // decode segments just like vue-router does
          .map(s => decodeURIComponent(s))));
    }
    return pageConfigNorm.value;
  });

  // We filter out special properties of the config when binding it to elements.
  // Without this we'd be trying to bind props like `children` which is a DOM property.
  const filterProps = {'layout': true, 'children': true};
  const filteredConfig = computed(() => Object.entries(configObj.value)
      .reduce((acc, [k, v]) => {
        if (Object.hasOwn(filterProps, k)) return acc;
        acc[k] = v;
        return acc;
      }, {}));

  return {
    pageConfig,
    pageConfigNorm,
    isLegacyOverview,
    isLegacySubPage,
    layout: computed(() => configObj.value.layout),
    config: filteredConfig
  };
}
