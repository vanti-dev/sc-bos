import {builtinWidgets} from '@/widgets/pallet.js';
import {computed, markRaw, reactive, toValue} from 'vue';
import {get as _get} from 'lodash';

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

  const {rootWidth, rootHeight} = svgRootSize(el);
  const bounds = elementBounds(el);
  const percent = (v, t) => {
    return (v / t * 100).toFixed(4) + '%';
  };

  el.style.visibility = 'hidden';

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
    bounds: {
      top: percent(bounds.y, rootHeight),
      left: percent(bounds.x, rootWidth),
      width: percent(bounds.width, rootWidth),
      height: percent(bounds.height, rootHeight)
    },
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

const svgRootSize = (el) => {
  const svg = el.ownerSVGElement;
  return {rootWidth: svg.width.baseVal.value, rootHeight: svg.height.baseVal.value};
};

const elementBounds = (el) => {
  const bBox = el.getBBox({stroke: true, markers: true});
  const ctm = el.getCTM();

  // note: browsers don't currently support {stroke: true} for getBBox so we fake it
  const style = window.getComputedStyle(el);
  const strokeWidth = style.getPropertyValue('stroke-width');
  if (strokeWidth) {
    const stroke = parseFloat(strokeWidth);
    const halfStroke = stroke / 2;
    bBox.x -= halfStroke;
    bBox.y -= halfStroke;
    bBox.width += stroke;
    bBox.height += stroke;
  }

  return matrixTransformRect(bBox, ctm);
};

/**
 * @param {DOMRect} rect
 * @param {DOMMatrix} matrix
 * @return {DOMRect}
 */
const matrixTransformRect = (rect, matrix) => {
  const tl = new DOMPoint(rect.x, rect.y).matrixTransform(matrix);
  const tr = new DOMPoint(rect.x + rect.width, rect.y).matrixTransform(matrix);
  const bl = new DOMPoint(rect.x, rect.y + rect.height).matrixTransform(matrix);
  const br = new DOMPoint(rect.x + rect.width, rect.y + rect.height).matrixTransform(matrix);

  const minx = Math.min(tl.x, tr.x, bl.x, br.x);
  const miny = Math.min(tl.y, tr.y, bl.y, br.y);
  const maxx = Math.max(tl.x, tr.x, bl.x, br.x);
  const maxy = Math.max(tl.y, tr.y, bl.y, br.y);

  return new DOMRect(minx, miny, maxx - minx, maxy - miny);
};
