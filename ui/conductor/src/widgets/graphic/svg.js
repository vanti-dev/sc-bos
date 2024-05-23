import {toValue} from '@/util/vue.js';
import binarySearch from 'binary-search';
import Color from 'colorjs.io';
import {get as _get} from 'lodash';
import {watch} from 'vue';

/**
 * Apply effects to the SVG based on the config and source data.
 */
export function useSvgEffects(el, config, sources) {
  for (const effect of effects) {
    if (!testEffect(effect, config)) continue;
    effect.apply(el, config, sources);
  }
}

const testEffect = (effect, cfg) => {
  if (typeof effect.test === 'string') return Boolean(cfg[effect.test]);
  if (typeof effect.test === 'function') return effect.test(cfg);
  return false;
};

const effects = [
  {
    test: 'fill',
    apply: (el, elementCfg, sources) => applyStyleColor('fill', el, elementCfg, sources)
  },
  {
    test: 'stroke',
    apply: (el, elementCfg, sources) => applyStyleColor('stroke', el, elementCfg, sources)
  }
];

const applyStyleColor = (prop, el, elementCfg, sources) => {
  const cfg = elementCfg[prop];
  const sourceCfg = cfg.source;
  const sourceResource = sources[sourceCfg.ref];
  if (!sourceResource) return;
  if (cfg.interpolate) {
    doColorInterpolation(
        () => _get(toValue(sourceResource.value), sourceCfg.property),
        cfg.interpolate.steps,
        color => el.style[prop] = color);
  }
};

const doColorInterpolation = (val, steps, onChange) => {
  const _steps = steps.map(s => {
    return {
      ...s,
      _color: new Color(s.color)
    };
  });
  watch(() => toValue(val), (value, oldValue) => {
    if (value === oldValue) return;
    const foundStep = binarySearch(_steps, value, (s, v) => s.value - v);
    if (foundStep >= 0) {
      // use exact color from the step
      onChange(_steps[foundStep]._color.toString());
    } else {
      // interpolate between steps
      const idx = ~foundStep;
      if (idx === 0) {
        onChange(_steps[0]._color.toString());
      } else if (idx === _steps.length) {
        onChange(_steps[_steps.length - 1]._color.toString());
      } else {
        const step1 = _steps[idx - 1];
        const step2 = _steps[idx];
        const ratio = (value - step1.value) / (step2.value - step1.value);
        onChange(step1._color.mix(step2._color, ratio).toString());
      }
    }
  }, {immediate: true});
};
