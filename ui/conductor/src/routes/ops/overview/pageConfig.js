import {builtinLayouts} from '@/layout/pallet.js';
import useBuildingConfig from '@/routes/ops/overview/pages/buildingConfig.js';
import useDashPage from '@/routes/ops/overview/pages/dashPage.js';
import {useUiConfigStore} from '@/stores/ui-config.js';
import {findActiveItem} from '@/util/router.js';
import {builtinWidgets} from '@/widgets/pallet.js';
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

  return {
    pageConfig,
    pageConfigNorm,
    isLegacyOverview,
    isLegacySubPage,
    layout: computed(() => configObj.value.layout),
    config: configObj
  };
}
