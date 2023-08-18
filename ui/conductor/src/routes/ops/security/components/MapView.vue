<template>
  <v-container fluid class="mb-0 mt-0 pb-0 pt-0 floor-plan__container">
    <PinchZoom @click="handleClick">
      <template #default="{ scale }">
        <Stack ref="groupingContainer">
          <!-- eslint-disable vue/no-v-html -->
          <div v-html="activeFloorPlan" ref="floorPlanSVG" :style="{ '--map-scale': scale }"/>
          <!-- eslint-enable vue/no-v-html -->
          <div v-if="showMenu" style="pointer-events: none">
            <div :style="calculateAnchorStyle" style="pointer-events: none">
              <HotPoint
                  v-slot="{ live }"
                  :item-key="elementWithMenu?.device?.name"
                  style="position: relative; top: 100%; transform-origin: 0 0; pointer-events: auto"
                  :style="{
                    transform: `scale(${1 / scale})`,
                  }">
                <AccessPointCard :device="elementWithMenu?.device" :paused="!live" @onClose="closeMenu"/>
              </HotPoint>
            </div>
          </div>
        </Stack>
      </template>
    </PinchZoom>
  </v-container>
</template>

<script setup>
import HotPoint from '@/components/HotPoint.vue';

import {useStatusBarStore} from '@/routes/ops/security/components/access-point-card/statusBarStore';
import Stack from '@/routes/ops/security/components/Stack.vue';
import PinchZoom from '@/routes/ops/security/map/PinchZoom.vue';
import {useAppConfigStore} from '@/stores/app-config';
import {convertSVGToPercentage} from '@/util/svg';
import {storeToRefs} from 'pinia';
import {computed, onBeforeUnmount, onMounted, reactive, ref, watch} from 'vue';
import AccessPointCard from './AccessPointCard.vue';

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
const elementWithMenu = reactive({
  device: null,
  source: null,
  target: null,
  x: 0,
  y: 0
});

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
  showClose.value = false;
  elementWithMenu.device = null;
  elementWithMenu.source = null;
  elementWithMenu.target = null;
  elementWithMenu.x = 0;
  elementWithMenu.y = 0;
};

/**
 * Find the device name in the props.deviceNames array
 *
 * @param {string} needle
 * @return {{name: string, source: string} | undefined}
 */
const findDevice = (needle) => {
  return props.deviceNames.find((deviceName) => deviceName.name.toLowerCase().endsWith('/' + needle.toLowerCase()));
};

/**
 *
 * @param {PointerEvent} event
 */
const handleClick = (event) => {
  const clickedElement = event.target.closest('[id]');

  // Find the parent group of the clicked element
  const parentGroup = clickedElement.closest('g[id^="doors_"]');
  if (!parentGroup) {
    return;
  }

  // Does the device sends a signal?
  const device = findDevice(clickedElement.id);
  if (!device) {
    return; // Do not proceed if the device does not send a signal
  }

  // Collect all the ids of the elements in the svg
  elementWithMenu.target = clickedElement;
  elementWithMenu.device = device;
  elementWithMenu.source = device.source;

  // Show menu
  showClose.value = true;
  showMenu.value = true;
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
        closeMenu();
        fetchFloorPlan(newValue).then((response) => {
          response.text().then((text) => {
            activeFloorPlan.value = text;
          });
        });
      }
    },
    {immediate: true, deep: true, flush: 'sync'}
);

// Watch for changes in the showClose prop then close menu
watch(
    showClose,
    (newValue, oldValue) => {
      if (newValue === false) {
        closeMenu();
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

.floor-plan__container ::v-deep path[id],
.floor-plan__container ::v-deep rect[id] {
    cursor: pointer;
}
</style>
