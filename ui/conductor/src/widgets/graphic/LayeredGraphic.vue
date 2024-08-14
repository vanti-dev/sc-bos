<template>
  <v-card :height="props.height" class="overflow-hidden pa-3 d-flex flex-column">
    <v-card-title v-if="props.title" class="text-h4 mb-0">{{ props.title }}</v-card-title>
    <overlay-stack v-if="props.fixed">
      <img v-if="bgSrc" :src="bgSrc" alt="Background for other layers">
      <graphic-layer
          v-for="(layer, i) in props.layers"
          :key="layer.title ?? i"
          :layer="layer"
          @click:element="onElementClick(layer, $event.element, $event.event)"
          :selected="selectionsByLayer[i]"
          @update:selected="onLayerSelectUpdate(i, $event)"/>
    </overlay-stack>
    <pinch-zoom v-else class="fill-height" :hide-controls="props.hideControls">
      <template #default>
        <overlay-stack>
          <img v-if="bgSrc" :src="bgSrc" alt="Background for other layers">
          <graphic-layer
              v-for="(layer, i) in props.layers"
              :key="layer.title ?? i"
              :layer="layer"
              @click:element="onElementClick(layer, $event.element, $event.event)"
              :selected="selectionsByLayer[i]"
              @update:selected="onLayerSelectUpdate(i, $event)"/>
        </overlay-stack>
      </template>
    </pinch-zoom>
  </v-card>
</template>

<script setup>
import {getMetadata} from '@/api/sc/traits/metadata.js';
import OverlayStack from '@/components/zoom/OverlayStack.vue';
import PinchZoom from '@/components/zoom/PinchZoom.vue';
import DeviceSideBar from '@/routes/devices/components/DeviceSideBar.vue';
import {useSidebarStore} from '@/stores/sidebar';
import GraphicLayer from '@/widgets/graphic/GraphicLayer.vue';
import {usePathUtils} from '@/widgets/graphic/path.js';
import {nameFromRequest} from '@/widgets/graphic/traits.js';
import {computed, ref, watch} from 'vue';

const props = defineProps({
  title: {
    type: String,
    default: null
  },
  height: {
    type: [Number, String],
    default: undefined
  },
  hideControls: {
    type: Boolean,
    default: false
  },
  fixed: {
    type: Boolean,
    default: false
  },
  background: {
    type: Object,
    default: null
  },
  layers: {
    type: Array,
    default: () => []
  }
});

const {toPath} = usePathUtils();
const bgSrc = computed(() => toPath(props.background?.svgPath));

const sidebar = useSidebarStore();
const selectionCtx = ref(null);
const onElementClick = async (layer, element, event) => {
  selectionCtx.value = {layer, element};
  // Find the name of the device we should be showing in the sidebar.
  // First we check if it's configured explicitly via the sidebar property.
  // Then we try to find a source that mentions a device name in the request.
  let name = element.sidebar?.name;
  if (!name) {
    for (const [, source] of Object.entries(element.sources ?? {})) {
      name = nameFromRequest(source.request);
      if (name) {
        break;
      }
    }
  }
  if (!name) return; // no name, nothing to show.

  // this is annoying, because the sidebar currently doesn't support fetching its own device info, so get it instead.
  const md = await getMetadata({name: name});
  if (md) {
    // open the sidebar with the metadata
    sidebar.title = md.appearance?.title ?? name;
    sidebar.data = {
      metadata: md, name,
      // these aren't used by the sidebar, but are used to work out if our selection is still the active one
      layer, element
    };
    sidebar.component = DeviceSideBar; // this line must be after the .data one!
    sidebar.visible = true;
  }
};

// force a single selection between all layers, and make sure we clear selection when someone else uses the sidebar
const selectionsByLayer = ref([]);
watch(() => props.layers, () => {
  selectionsByLayer.value = props.layers.map(() => null);
});
watch(() => sidebar.data, (c) => {
  if (c.element !== selectionCtx.value?.element || c.layer !== selectionCtx.value?.layer) {
    // someone else selected something, or the sidebar was hidden
    selectionsByLayer.value = selectionsByLayer.value.map(() => null);
    selectionCtx.value = null;
  }
});
const onLayerSelectUpdate = (layerIdx, selected) => {
  for (let i = 0; i < props.layers.length; i++) {
    if (i === layerIdx) selectionsByLayer.value[i] = selected;
    else selectionsByLayer.value[i] = null;
  }
};
</script>

<style scoped>

</style>
