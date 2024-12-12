<template>
  <v-container fluid class="mb-0 mt-0 pt-0 floor-plan__container">
    <pinch-zoom @click="handleClick">
      <template #default="{ scale }">
        <overlay-stack ref="groupingContainer">
          <!-- eslint-disable-next-line vue/no-v-html -->
          <div v-html="activeFloorPlan" ref="floorPlanSVG" :style="{ '--map-scale': scale }"/>
          <div v-if="showMenu" style="pointer-events: none">
            <div :style="calculateAnchorStyle" style="pointer-events: none">
              <hot-point
                  v-slot="{ live }"
                  :item-key="elementWithMenu?.device?.name"
                  style="position: relative; top: 100%; transform-origin: 0 0; pointer-events: auto"
                  :style="{
                    transform: `scale(${1 / scale})`,
                  }">
                <access-point-card
                    :device="elementWithMenu?.device"
                    :paused="!live"
                    @click:close="closeMenu"
                    show-close/>
              </hot-point>
            </div>
          </div>
          <div class="door-status-tracker">
            <hot-point
                v-slot="{ live }"
                v-for="door in doors"
                :key="door.name"
                :item-key="door.name"
                class="door-status-tracker__item"
                :style="door.style">
              <!-- If door has Access data reading and has no OpenClose reading -->
              <with-access
                  v-if="hasTrait(door.name, 'Access') && !hasTrait(door.name, 'OpenClose')"
                  :name="door.name"
                  :paused="!live"
                  v-slot="{ resource: accessResource }">
                <with-status :name="door.name" :paused="!live" v-slot="{ resource: statusResource }">
                  <door-color
                      :name="door.name"
                      :access-attempt="accessResource.value"
                      :status-log="statusResource.value"
                      class="door-status-tracker__item"
                      @updateFill="setDoorFill"/>
                </with-status>
              </with-access>
              <!-- If door has no Access data reading and has OpenClose reading -->
              <with-open-close
                  v-if="!hasTrait(door.name, 'Access') && hasTrait(door.name, 'OpenClose')"
                  :name="door.name"
                  :paused="!live"
                  v-slot="{ resource: openCloseResource }">
                <with-status :name="door.name" :paused="!live" v-slot="{ resource: statusResource }">
                  <door-color
                      :name="door.name"
                      :open-close="openCloseResource.value"
                      :status-log="statusResource.value"
                      class="door-status-tracker__item"
                      @updateStroke="setDoorStroke"/>
                </with-status>
              </with-open-close>
              <!-- If door has Access data reading and has OpenClose reading -->
              <with-access
                  v-if="hasTrait(door.name, 'Access') && hasTrait(door.name, 'OpenClose')"
                  :name="door.name"
                  :paused="!live"
                  v-slot="{ resource: accessResource }">
                <with-open-close :name="door.name" :paused="!live" v-slot="{ resource: openCloseResource }">
                  <with-status :name="door.name" :paused="!live" v-slot="{ resource: statusResource }">
                    <door-color
                        :name="door.name"
                        :access-attempt="accessResource.value"
                        :open-close="openCloseResource.value"
                        :status-log="statusResource.value"
                        class="door-status-tracker__item"
                        @updateFill="setDoorFill"
                        @updateStroke="setDoorStroke"/>
                  </with-status>
                </with-open-close>
              </with-access>
            </hot-point>
          </div>
        </overlay-stack>
      </template>
    </pinch-zoom>
  </v-container>
</template>

<script setup>
import HotPoint from '@/components/HotPoint.vue';
import OverlayStack from '@/components/zoom/OverlayStack.vue';
import PinchZoom from '@/components/zoom/PinchZoom.vue';
import DoorColor from '@/routes/ops/security/components/DoorColor.vue';
import {useUiConfigStore} from '@/stores/uiConfig.js';
import WithAccess from '@/traits/access/WithAccess.vue';
import WithOpenClose from '@/traits/openClose/WithOpenClose.vue';
import WithStatus from '@/traits/status/WithStatus.vue';
import {subPath} from '@/util/path.js';
import {convertSVGToPercentage} from '@/util/svg';
import {computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch} from 'vue';
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
const uiConfig = useUiConfigStore();
const showClose = ref(false);
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
const svgVersion = ref(0);
const getSvgEl = () => {
  svgVersion.value; // register interest in this ref
  const svgContainer = floorPlanSVG.value;
  if (!svgContainer) {
    return undefined;
  }
  const svgEl = svgContainer?.querySelector('svg');
  if (!svgEl) {
    return undefined;
  }
  return svgEl;
};
const getSVGViewBox = () => {
  const svgEl = getSvgEl();
  if (!svgEl) {
    return undefined;
  }
  const [x, y, w, h] = svgEl.getAttribute('viewBox').split(' ');

  return {
    x: parseInt(x),
    y: parseInt(y),
    width: parseInt(w),
    height: parseInt(h)
  };
};

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
  const viewBox = getSVGViewBox();
  if (!viewBox) {
    return {};
  }

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
const floorPlanSVGPath = computed(() => {
  const floorName = props.floor;
  const svgPath = subPath(
      uiConfig.config.siteFloorPlans.find((floorPlan) => floorPlan.name === floorName)?.svgPath,
      uiConfig.configUrl
  );
  return svgPath;
});
/**
 * Fetch function to get the floor plan svg
 *
 * @param {string} svgPath
 * @return {Promise<Response>}
 */
const fetchFloorPlan = async (svgPath) => {
  // Fetch the floor plan svg
  // Don't forget to add ?raw to the end of the url to get the raw svg (string injected into v-html)
  const response = await fetch(svgPath, {
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
  if (!clickedElement) {
    return;
  }

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

// doors contains door names, location on the map, and the svg element.
const doors = computed(() => {
  if (!floorPlanSVG.value) {
    return [];
  }

  const viewBox = getSVGViewBox();
  if (!viewBox) {
    return [];
  }
  const svgEl = getSvgEl();
  if (!svgEl) {
    return [];
  }
  return props.deviceNames
      .map((device) => {
        const localId = device.name.split('/').pop();
        const door = svgEl.querySelector(`#${localId}`);
        if (!door) return undefined;

        const elRect = door.getBBox();
        const percentage = convertSVGToPercentage(viewBox, elRect);

        const inset = 0.02;
        const x = (percentage.x - inset) * 100;
        const y = (percentage.y - inset) * 100;
        const width = (percentage.width + 2 * inset) * 100;
        const height = (percentage.height + 2 * inset) * 100;
        return {
          name: device.name,
          el: door,
          style: {
            width: `${width}%`,
            height: `${height}%`,
            left: `${x}%`,
            top: `${y + 1}%`
          }
        };
      })
      .filter((d) => Boolean(d));
});
// doorColors contains the intended colour for each door.
// We keep this as a data structure instead of just setting the value in case we know what colour it should be
// before the svg is loaded.
const doorFills = ref({});
const doorStrokes = ref({});
const setDoorFill = ({name, color}) => {
  doorFills.value[name] = color;
};
const setDoorStroke = ({name, color}) => {
  doorStrokes.value[name] = color;
};

// watch for changes in the colours and svg and invoke dom actions to update the svg.
watch(
    [doorFills, doorStrokes, doors],
    () => {
      doors.value.forEach(({el, name}) => {
        const color = doorFills.value[name] ?? 'grey';
        const stroke = doorStrokes.value[name] ?? 'unknown';

        if (el) {
          el.removeAttribute('style');

          // Handle fill class updates
          el.classList.remove('success', 'error', 'warning', 'grant_unknown', 'grey');
          el.classList.add(color);

          // Handle stroke class updates
          el.classList.remove('open', 'closed', 'moving', 'unknown');
          el.classList.add(stroke);

          const fill = doorFills.value[name] ? doorFills.value[name] : doorStrokes.value[name] + '-fill';
          el.classList.remove('open-fill', 'closed-fill', 'moving-fill', 'unknown-fill');
          el.classList.add(fill);
        }
      });
    },
    {deep: true}
);


const hasTrait = (device, traitName) => {
  const traits = {};
  let traitFullName;

  const findDevice = props.deviceNames.find((deviceName) => {
    return deviceName.name === device;
  });

  if (!findDevice) return false;

  if (traitName === 'OpenClose') traitFullName = 'smartcore.traits.OpenClose';
  else if (traitName === 'Access') traitFullName = 'smartcore.bos.Access';

  if (findDevice.traits.includes(traitFullName)) {
    traits[traitName] = true;
  }

  return traits[traitName];
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
    floorPlanSVGPath,
    (newValue, oldValue) => {
      if (newValue !== oldValue) {
        closeMenu();
        fetchFloorPlan(newValue).then((response) => {
          response.text().then((text) => {
            activeFloorPlan.value = text;
            nextTick(() => {
              svgVersion.value++;
            });
          });
        });
      }
    },
    {immediate: true, deep: true, flush: 'sync'}
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

.floor-plan__container :deep(path[id]),
.floor-plan__container :deep(rect[id]) {
  cursor: pointer;
}

.door-status-tracker {
  position: relative;
  pointer-events: none;
}

.door-status-tracker__item {
  position: absolute;
}

:deep(svg .success) {
  fill: rgb(var(--v-theme-success));
}

:deep(svg .warning) {
  fill: rgb(var(--v-theme-warning));
}

:deep(svg .error) {
  fill: rgb(var(--v-theme-error));
}

:deep(svg .open),
:deep(svg .moving) {
  stroke: rgb(var(--v-theme-warning));
  stroke-width: 125px;
  transition: all 0.5s ease-in-out;
}

:deep(svg .closed) {
  stroke: rgb(var(--v-theme-success));
  stroke-width: 75px;
  transition: all 0.5s ease-in-out;
}

:deep(svg .unknown) {
  stroke: #ffffff5e;
  stroke-width: 75px;
  transition: all 0.5s ease-in-out;
}

:deep(svg .open-fill),
:deep(svg .moving-fill) {
  fill: rgb(var(--v-theme-warning));
  transition: all 0.5s ease-in-out;
}

:deep(svg .closed-fill) {
  fill: rgb(var(--v-theme-success));
  transition: all 0.5s ease-in-out;
}

:deep(svg .unknown-fill) {
  fill: #ffffff85;
  transition: all 0.5s ease-in-out;
}
</style>
