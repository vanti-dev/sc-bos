<template>
  <v-container fluid class="mb-0 mt-0 pb-0 pt-0 floor-plan__container">
    <PinchZoom @click="handleClick">
      <template #default="{ scale }">
        <Stack ref="groupingContainer">
          <!-- eslint-disable vue/no-v-html -->
          <div v-html="activeFloorPlan" ref="floorPlanSVG" :style="{ '--map-scale': scale }"/>
          <!-- eslint-enable vue/no-v-html -->
          <div v-if="showMenu">
            <div :style="calculateAnchorStyle">
              <WithAccess
                  v-slot="{ resource }"
                  :name="elementWithMenu?.deviceId"
                  style="position: relative; top: 100%; transform-origin: 0 0"
                  :style="`transform: scale(${1 / scale})`">
                <AccessPointCard
                    v-bind="resource"
                    :name="elementWithMenu?.deviceId"
                    show-close
                    :source="elementWithMenu?.source"
                    @onClose="closeMenu"/>
              </WithAccess>
            </div>
          </div>
        </Stack>
      </template>
    </PinchZoom>
  </v-container>
</template>

<script setup>
import {computed, onMounted, onBeforeUnmount, reactive, ref, watch} from 'vue';
import {useAppConfigStore} from '@/stores/app-config';
import PinchZoom from '@/routes/ops/security/map/PinchZoom.vue';
import Stack from '@/routes/ops/security/components/Stack.vue';
import WithAccess from '@/routes/devices/components/renderless/WithAccess.vue';
import AccessPointCard from './AccessPointCard.vue';

import {convertSVGToPercentage} from '@/util/svg';

// -------------- Props -------------- //
const props = defineProps({
  deviceNames: {
    type: Array,
    default: () => []
  },
  floor: {
    type: String,
    default: 'level0' // TODO: change to actual ground floor
  }
});

// -------------- Data & Reactive References -------------- //
const {config} = useAppConfigStore();
const activeFloorPlan = ref('');
const floorPlanSVG = ref(null);
const groupingContainer = ref(null);

const showMenu = ref(false);
let elementWithMenu = reactive({
  deviceId: null,
  source: null,
  target: null,
  x: 0,
  y: 0
});
const groupedIds = ref({});


// -------------- Computed Properties -------------- //
const getSVGViewBox = computed(() => {
  const [x, y, w, h] = floorPlanSVG.value.querySelector('svg').getAttribute('viewBox').split(' ');

  return {
    x: parseInt(x),
    y: parseInt(y),
    width: parseInt(w),
    height: parseInt(h)
  };
});


const getClickedRectBBox = computed(() => {
  if (!elementWithMenu.target) {
    return {};
  }

  return elementWithMenu.target.getBBox();
});

const calculateAnchorStyle = computed(() => {
  if (!elementWithMenu.target || !groupingContainer.value) {
    return {};
  }

  // Get the bounding rectangle of the SVG element
  const clickedRect = getClickedRectBBox.value;
  const viewBox = getSVGViewBox.value;

  const percentage = convertSVGToPercentage(viewBox, clickedRect);

  const x = percentage.x * 100;
  const y = percentage.y * 100;
  const width = percentage.width * 100;
  const height = percentage.height * 100;

  return {
    width: `${width}%`,
    height: `${height}%`,
    left: `${x}%`,
    top: `${y + 1}%`,
    position: 'relative'
  };
});

// -------------- Methods -------------- //
/**
 * Fetch function to get the floor plan svg
 *
 * @param {string} selectedFloor
 * @return {Promise<Response>}
 */
const fetchFloorPlan = async (selectedFloor) => {
  const floorPlan = config.siteFloorPlans.find((floorPlan) => floorPlan.name === selectedFloor);

  // Fetch the floor plan svg
  // Don't forget to add ?raw to the end of the url to get the raw svg (string injected into v-html)
  const response = await fetch(floorPlan.svgPath + '?raw', {
    headers: {
      'Content-Type': 'image/svg+xml'
    }
  });
  return response;
};

// Close the menu with X button click
const closeMenu = () => {
  showMenu.value = false;
  elementWithMenu = {deviceId: null, source: null, target: null};
};

/**
 * Collecting all the ids of the elements in the svg
 * and grouping them by the parent group id
 *
 * @param {HTMLElement} element
 */
const traverseAndCollectIds = (element) => {
  // If the element is a 'g', then we simply ignore and don't traverse its children.
  if (element.tagName === 'g' || ['detail', 'outline'].includes(element.id)) {
    return;
  }

  // If the element is a <path> or <rect> and is not inside a <g>, then collect its ID.
  if (element.tagName === 'path' || element.tagName === 'rect') {
    const parts = element.id.split('_');
    if (parts.length > 1) {
      const groupName = parts[0] + 's';
      const itemId = parts.slice(1).join('_'); // This is to handle cases with multiple underscores

      if (!groupedIds.value[groupName]) {
        groupedIds.value[groupName] = [];
      }
      groupedIds.value[groupName].push(itemId);
    }
  }

  // Continue traversal for other children.
  const children = element.children;
  for (let i = 0; i < children.length; i++) {
    traverseAndCollectIds(children[i]);
  }
};

/**
 * Find the device name in the props.deviceNames array
 *
 * @param {HTMLElement} element
 * @return {{name: string, source: string} | undefined}
 */
const findDevice = (element) => {
  return props.deviceNames.find((deviceName) => deviceName.name === element.id.split('_')[1]);
};

/**
 *
 * @param {PointerEvent} event
 */
const handleClick = (event) => {
  // Do not react to clicks on the svg itself (blank space around the floor plan)
  if (event.target.tagName === 'svg') {
    return;
  }

  const clickedElement = event.target;

  // Does the device sends a signal?
  if (!findDevice(clickedElement)) {
    return; // Do not proceed if the device does not send a signal
  }

  // Collect all the ids of the elements in the svg
  elementWithMenu.target = clickedElement;
  elementWithMenu.deviceId = findDevice(clickedElement).name;
  elementWithMenu.source = findDevice(clickedElement).source;

  // Show menu
  showMenu.value = true;
};

/**
 * Add or remove event listeners and MutationObserver
 *
 * @param {string} action
 */
const manageEventListeners = (action) => {
  floorPlanSVG.value[action + 'EventListener']('click', handleClick);
};

// -------------- Lifecycle Hooks -------------- //
onMounted(() => {
  manageEventListeners('add');
});

onBeforeUnmount(() => {
  manageEventListeners('remove');
});

// -------------- Watchers -------------- //
// Watch for changes in the floor prop then
// fetch the floor plan svg
watch(
    () => props.floor,
    (newValue, oldValue) => {
      if (newValue !== oldValue) {
        floorPlanSVG.value = document.getElementById(newValue);
        fetchFloorPlan(newValue).then((response) => {
          response.text().then((text) => {
            activeFloorPlan.value = text;
          });
        });
      }
    },
    {immediate: true, deep: true, flush: 'sync'}
);

// Watch for floorPlanSVG changes and traverse the svg to collect all the ids
watch(
    floorPlanSVG,
    (newValue, oldValue) => {
      if (newValue && newValue !== oldValue) {
        traverseAndCollectIds(newValue);
      }
    },
    {immediate: true}
);
</script>

<style lang="scss" scoped>
.floor-plan__container {
  position: relative;
  /* fill the container, minus the top bar and sc status bar */
  height: calc(100vh - 215px);
  overflow: hidden;
}

.floor-plan__container .pinch-zoom {
  /* fill the container so that zoom controls show in the bottom-right */
  height: 100%;
}

.pinch-zoom {
  /* defaults, overridden in the template & deviceMarkers() */
  --map-scale: 1;
  --translate-x: 0;
  --translate-y: 0;
}
</style>
