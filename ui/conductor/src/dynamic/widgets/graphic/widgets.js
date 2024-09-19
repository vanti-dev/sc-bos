import {elementBounds, svgRootSize} from '@/util/svg.js';
import {builtinWidgets} from '@/dynamic/widgets/pallet.js';
import {get as _get} from 'lodash';
import {computed, markRaw, reactive, toValue} from 'vue';

/**
 * Calculates widget placement and props for a single source element and effect
 *
 * @param {SVGGraphicsElement} el - the target SVG element to associate the widget with
 * @param {layer.Element} config
 * @param {Record<string, RemoteResource<any>>} sources
 * @return {layer.WidgetInstance | null}
 */
export function useWidgetEffects(el, config, sources) {
  const effect =
      /** @type {layer.WidgetEffect | undefined} */
      config.effects.find(effect => effect.type === 'widget');
  if (!effect) return null;

  const comp = loadComponent(effect.component);
  if (!comp) {
    console.warn(`Unknown widget component: ${effect.component}`);
    return null;
  }

  // allow selecting a child element for the widget to replace
  if (effect.selector) {
    el = el.querySelector(effect.selector);
    if (!el) {
      console.warn(`Could not find child element with selector: ${effect.selector}`);
      return null;
    }
  }

  const {rootWidth, rootHeight} = svgRootSize(el);
  const bounds = elementBounds(el);
  const percent = (v, t) => {
    return (v / t * 100).toFixed(4) + '%';
  };
  const boundsPercent = {
    top: percent(bounds.y, rootHeight),
    left: percent(bounds.x, rootWidth),
    width: percent(bounds.width, rootWidth),
    height: percent(bounds.height, rootHeight)
  };

  if (!effect.showElement) {
    el.style.visibility = 'hidden';
  }

  const props = {};
  for (const [k, v] of Object.entries(effect.props)) {
    if (typeof v === 'object' && 'ref' in v) {
      const source = sources[v.ref];
      if (source) {
        props[k] = computed(() => _get(toValue(source.value), v.property));
      } else {
        console.warn(`Unknown source: ${v.ref}`);
      }
    } else {
      props[k] = v;
    }
  }

  return {
    component: markRaw(comp),
    bounds: boundsPercent,
    props: reactive(props)
  };
}

const loadComponent = (compString) => {
  if (compString.startsWith('builtin:')) {
    const [, builtin] = compString.split(':');
    return markRaw(builtinWidgets[builtin]);
  }
  return null;
};
