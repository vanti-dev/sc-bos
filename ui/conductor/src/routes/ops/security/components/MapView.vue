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
              <AccessPointCard
                  :device="findDevice(elementWithMenu?.deviceId)"
                  style="position: relative; top: 100%; transform-origin: 0 0"
                  :style="`transform: scale(${1 / scale})`"
                  @onClose="closeMenu"/>
            </div>
          </div>
        </Stack>
      </template>
    </PinchZoom>
  </v-container>
</template>

<script setup>
import {computed, onMounted, onBeforeUnmount, reactive, ref, watch} from 'vue';
import {storeToRefs} from 'pinia';
import {useAppConfigStore} from '@/stores/app-config';
import PinchZoom from '@/routes/ops/security/map/PinchZoom.vue';
import Stack from '@/routes/ops/security/components/Stack.vue';
import AccessPointCard from './AccessPointCard.vue';

import {useStatusBarStore} from '@/routes/ops/security/components/access-point-card/statusBarStore';
import {convertSVGToPercentage} from '@/util/svg';

// -------------- Props -------------- //
const props = defineProps({
  deviceNames: {
    type: Array,
    default: () => []
  },
  floor: {
    type: String,
    default: 'Ground Floor'
  }
});

// -------------- Data & Reactive References -------------- //
const {config} = useAppConfigStore();
const {showClose} = storeToRefs(useStatusBarStore());
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
  // Check if the element is a group with an ID containing 'door' or 'doors'
  if (element.tagName === 'g' && element.id && element.id.toLowerCase().includes('door')) {
    const groupKey = element.id.split('_')[1];

    // If the key doesn't exist in the dictionary, create an empty array for it
    if (!groupedIds.value[groupKey]) {
      groupedIds.value[groupKey] = [];
    }

    // Collecting the IDs of <path> and <rect> elements
    Array.from(element.children).forEach((child) => {
      if (['path', 'rect'].includes(child.tagName) && child.id) {
        groupedIds.value[groupKey].push(child.id.toLowerCase());
      }
    });
  }

  // Continue traversal for other children
  Array.from(element.children).forEach((child) => {
    traverseAndCollectIds(child);
  });
};

/**
 * Find the device name in the props.deviceNames array
 *
 * @param {HTMLElement} element
 * @return {{name: string, source: string} | undefined}
 */
const findDevice = (element) => {
  if (element.id) {
    return props.deviceNames.find((deviceName) => deviceName.name.toLowerCase() === element.id.toLowerCase());
  } else return props.deviceNames.find((deviceName) => deviceName.name.toLowerCase() === element.toLowerCase());
};

/**
 *
 * @param {PointerEvent} event
 */
const handleClick = (event) => {
  const clickedElement = event.target;

  // Do not react to clicks on the svg itself (blank space around the floor plan)
  if (clickedElement.tagName === 'svg') {
    return;
  }

  // Check if the parent of the clicked element is 'outline' or 'detail' group.
  const parentGroup = clickedElement.parentElement;
  if (parentGroup && (parentGroup.id === 'outline' || parentGroup.id === 'detail')) {
    return;
  }

  // Check if the parent group of the clicked element contains 'door' or 'doors'
  if (!parentGroup || !parentGroup.id || !parentGroup.id.toLowerCase().includes('door')) {
    return; // Do not proceed if the clicked element is not inside a 'door' or 'doors' group
  }

  // Does the device sends a signal?
  const device = findDevice(clickedElement);
  if (!device) {
    return; // Do not proceed if the device does not send a signal
  }

  // Collect all the ids of the elements in the svg
  elementWithMenu.target = clickedElement;
  elementWithMenu.deviceId = device.name;
  elementWithMenu.source = device.source;

  // Show menu
  showMenu.value = true;
  showClose.value = true;
};

/**
 * Add or remove event listeners
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

// Watch for changes in the showClose prop then close menu
watch(showClose, (newValue) => {
  if (!newValue) {
    closeMenu();
  }
});

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
