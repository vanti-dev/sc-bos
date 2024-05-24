<template>
  <!-- eslint-disable vue/no-v-html -->
  <div
      v-html="svgRaw"
      class="svg--container"
      @click="onSvgClick"
      ref="svgContainerEl"/>
</template>

<script setup>
import {closeResource} from '@/api/resource.js';
import {usePullTrait} from '@/traits/traits.js';
import {loadConfig} from '@/widgets/graphic/config.js';
import {usePathUtils} from '@/widgets/graphic/path.js';
import {useSvgEffects} from '@/widgets/graphic/svg.js';
import {effectScope, nextTick, onUnmounted, reactive, ref, watch} from 'vue';

const props = defineProps({
  layer: {
    type: Object,
    required: true
  },
  // Either null (for no selection), the index of the element that is selected, or an array of selected element indexes
  selected: {
    type: [Number, Array],
    default: null
  },
  // control whether one or multiple elements can be selected
  multiple: {
    type: Boolean,
    default: false
  }
});
const emit = defineEmits(['click:element', 'update:selected']);

const _selected = ref(null);
watch(() => props.selected, (v) => _selected.value = v);

/**
 * @param {PointerEvent} event
 */
const onSvgClick = (event) => {
  const el = event.target.closest('[data-element-idx]');
  if (!el) return; // not an interesting click
  const elementIdx = parseInt(el.dataset.elementIdx);
  const layerElement = config.value.elements[elementIdx];
  emit('click:element', {element: layerElement, index: elementIdx, event});
  selectElement(el, elementIdx, layerElement);
};

const selectElement = (el, elementIdx, element) => {
  if (props.multiple) {
    const oldIdx = _selected.value.indexOf(elementIdx);
    if (oldIdx === -1) {
      _selected.value.push(elementIdx);
    } else {
      _selected.value.splice(oldIdx, 1);
    }
  } else {
    if (_selected.value === elementIdx) {
      _selected.value = null;
    } else {
      _selected.value = elementIdx;
    }
  }
  emit('update:selected', _selected.value);
};
// Convert the selected value from singular to multiple as needed.
// Don't send events as the parent changed this.
watch(() => props.multiple, (multiple, oldMultiple) => {
  if (multiple === oldMultiple) return;
  if (multiple) {
    if (_selected.value === null) _selected.value = [];
    else _selected.value = [_selected.value];
  } else {
    if (_selected.value?.length > 0) _selected.value = _selected.value[0];
    else _selected.value = null;
  }
});

// track selected elements and update dom highlighting, by changing a class
const normSelected = (s) => {
  if (s === null || s === undefined) return new Set();
  if (Array.isArray(s)) return new Set(s);
  return new Set([s]);
};
watch(_selected, (selected, oldSelected) => {
  const nowSelected = normSelected(selected);
  const wasSelected = normSelected(oldSelected);

  // remove 'selected' class from everything that was selected, but isn't now
  for (const was of wasSelected) {
    if (nowSelected.has(was)) continue;
    const els = svgEl.value.querySelectorAll(`[data-element-idx="${was}"]`);
    for (const el of els) {
      el.classList.remove('selected');
    }
  }
  // add 'selected' class to everything that is selected now, but wasn't before
  for (const now of nowSelected) {
    if (wasSelected.has(now)) continue;
    const els = svgEl.value.querySelectorAll(`[data-element-idx="${now}"]`);
    for (const el of els) {
      el.classList.add('selected');
    }
  }
}, {immediate: true, deep: true});

// fetch data from the server based on the layer configuration
const svgRaw = ref('');
const config = ref(null);

const {toPath} = usePathUtils();

const fetchLayer = async (layer) => {
  const gotConfig = loadConfig(toPath(layer.configPath));
  const gotSvg = fetch(toPath(layer.svgPath))
      .then((res) => res.text());
  return {
    config: await gotConfig,
    svgRaw: await gotSvg
  };
};
watch(() => props.layer, async () => {
  if (!props.layer) return;
  const res = await fetchLayer(props.layer);
  svgRaw.value = res.svgRaw;
  config.value = res.config;
}, {immediate: true});

// adjust/inspect the dom to make future interactions easier/faster
const svgContainerEl = ref(null);
const svgEl = ref(null);
watch([svgContainerEl, svgRaw], ([containerEl, svgRaw]) => {
  if (!containerEl || !svgRaw) return;
  nextTick(() => {
    svgEl.value = containerEl.querySelector('svg');
    if (!svgEl.value) {
      console.warn('no svg element found as child of svgContainer', containerEl);
    }
  });
});

const annotateSvgDom = (svgEl, config) => {
  for (let i = 0; i < config.elements.length; i++) {
    const le = config.elements[i];
    const els = svgEl.querySelectorAll(le.selector);
    if (!els) {
      console.warn('layer element not found for selector', le.selector, props.layer.title);
      continue;
    }
    // we use this attribute during clicks to find
    // a. the correct element to use as the click target
    // b. to find the layer element that describes what to do with the click
    for (const el of els) {
      el.setAttribute('data-element-idx', '' + i);
    }
  }
};
watch([svgEl, config], ([svgEl, config]) => {
  if (!svgEl || !config) return;
  annotateSvgDom(svgEl, config);
});

// setup any effects on the svg elements based on the sources of data

// A cache of all the server request resources.
const scopeClosers = ref(/** @type {(function():void)[]} */ []);
const closeAll = () => {
  scopeClosers.value.forEach(r => r());
  scopeClosers.value = [];
};
onUnmounted(() => {
  closeAll();
});
watch([svgEl, config], ([svgEl, config]) => {
  closeAll();
  if (!svgEl || !config) return; // do nothing, not ready yet

  const scope = effectScope();
  scope.run(() => {
    for (const element of config.elements ?? []) {
      if (!element.sources) continue; // no source of info, so skip
      const els = svgEl.querySelectorAll(element.selector);
      if (!els) continue; // no point continuing as we can't update the element

      // capture information from the server
      const sources = {};
      for (const [name, source] of Object.entries(element.sources)) {
        const resource = usePullTrait(source.trait, source.request);
        scopeClosers.value.push(() => closeResource(reactive(resource)));
        sources[name] = resource;
      }

      // setup dom changes based on server collected data
      for (const el of els) {
        useSvgEffects(el, element, sources);
      }
    }
  });
  scopeClosers.value.push(() => scope.stop());
});
</script>

<style scoped>
.svg--container {
  display: grid;
  align-items: stretch;
  justify-items: stretch;
  /** make sure clicks pass through any area that isn't marked for interaction */
  pointer-events: none;
}

.svg--container > ::v-deep(svg) {
  /** fix svgs that have width/height attributes in them **/
  width: auto !important;
  height: auto !important;
}

.svg--container ::v-deep([data-element-idx]) {
  cursor: pointer;
  pointer-events: auto;
  transition: filter 0.2s cubic-bezier(.25, .8, .25, 1);
}

.svg--container ::v-deep([data-element-idx].selected) {
  filter: drop-shadow(0 6px 10px rgba(0, 0, 0, 0.19)) drop-shadow(0 3px 6px rgba(0, 0, 0, 0.63));
}
</style>
