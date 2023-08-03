<template>
  <v-container fluid class="mb-0 pb-0 floor-plan__container">
    <PinchZoom @click="handleClick(event)">
      <template #default="{ scale }">
        <!-- eslint-disable-next-line vue/no-v-html -->
        <div v-html="activeFloorPlan" :id="props.floor" ref="svgContainer" :style="{ '--map-scale': scale }"/>
        <div>
          <v-menu
              v-model="showMenu"
              absolute
              bottom
              :position-x="elementWithMenu.x"
              :position-y="elementWithMenu.y"
              transition="fade-transition-v2">
            <template #activator="{ on }">
              <div v-bind="on" ref="menuActivator"/>
            </template>
            <v-tooltip bottom>
              <template #activator="{ on, attrs }">
                <v-btn
                    class="mt-2 elevation-4"
                    color="grey darken-2"
                    fab
                    x-small
                    style="top: -5px; left: 91.5%"
                    v-bind="attrs"
                    v-on="on"
                    @click="closeMenu">
                  <v-icon>mdi-close</v-icon>
                </v-btn>
              </template>
              <span>Close</span>
            </v-tooltip>
            <AccessPointCard/>
          </v-menu>
        </div>
      </template>
    </PinchZoom>
  </v-container>
</template>

<script setup>
import {onMounted, onUnmounted, ref, watch} from 'vue';
import {useAppConfigStore} from '@/stores/app-config';

import PinchZoom from '@/routes/ops/security/map/PinchZoom.vue';
import AccessPointCard from './AccessPointCard.vue';

// ------------------- Props ------------------- //
const props = defineProps({
  floor: {
    type: String,
    default: 'level0'
  }
});

// ------------------- Data ------------------- //
const {config} = useAppConfigStore();

const activeFloorPlan = ref('');
const floorPlanSVG = ref(null);
const elementWithMenu = ref({
  target: '' | null,
  x: 0,
  y: 0
});
const showMenu = ref(false);
const closeMenu = () => {
  showMenu.value = false;
  elementWithMenu.value = {
    target: null,
    x: 0,
    y: 0
  };
};
const groupedIds = ref({});

// ------------------- Methods ------------------- //

/**
 * Collecting all the ids of the elements in the svg
 * and grouping them by the parent group id
 *
 * @param {HTMLElement} element
 */
const traverseAndCollectIds = (element) => {
  const svgContainer = element.id === props.floor;
  const svgGroup = element.tagName === 'g';
  let group;

  if (!svgContainer && element.id) {
    group = svgGroup ? element.id : element.closest('g')?.id;
    if (group) {
      if (!groupedIds.value[group]) {
        groupedIds.value[group] = [];
      }
    }
  }

  if (!svgContainer && !svgGroup && element.id && group === element.closest('g')?.id) {
    groupedIds.value[group].push(element.id);
  }

  const children = element.children;
  for (let i = 0; i < children.length; i++) {
    traverseAndCollectIds(children[i]);
  }
};

/**
 *
 * @param {HTMLElement} element
 * @return {HTMLElement|undefined}
 */
function findDeepestChild(element) {
  let deepestChild = element;

  while (deepestChild.lastElementChild) {
    deepestChild = deepestChild.lastElementChild;
  }

  if (deepestChild.id) {
    return deepestChild;
  } else return;
}

/**
 *
 * @param {MouseEvent} event
 */
function handleClick(event) {
  // Find deepest child
  const clickedElement = event.target;
  elementWithMenu.value.target = findDeepestChild(clickedElement);
  // elementWithMenu.value.target = findDeepestChild(clickedElement);

  // If no child found, return
  if (!elementWithMenu.value.target) {
    return;
  }

  // Reset menu state
  elementWithMenu.value.target = null;
  elementWithMenu.value.x = 0;
  elementWithMenu.value.y = 0;

  // Calculate menu position
  const clickedRect = clickedElement.getBoundingClientRect();

  // calculate position to the left
  elementWithMenu.value.x = clickedRect.left - Math.floor((clickedRect.left / 100) * 25);
  // calculate position to the bottom
  elementWithMenu.value.y = clickedRect.top + clickedRect.height * 1.5;

  // Show/Hide menu
  showMenu.value = true;
}

// -------------------- //
// fetch function to get the floor plan svg
const fetchFloorPlan = async () => {
  // TODO: change level0 to props.floor
  const floorPlan = config.siteFloorPlans.find((floorPlan) => floorPlan.name === 'level0');

  // fetch the svg
  // add ?raw to the end of the url to get the raw svg
  const response = await fetch(floorPlan.svgPath + '?raw', {
    headers: {
      'Content-Type': 'image/svg+xml'
    }
  });
  return response;
};

// Watch for changes in the floor prop then
// fetch the floor plan svg
watch(
    () => props.floor,
    (newValue, oldValue) => {
      if (newValue !== oldValue) {
        fetchFloorPlan().then((response) => {
          response.text().then((text) => {
            activeFloorPlan.value = text;
          });
        });
      }
    },
    {immediate: true}
);

// ------------------- Lifecycle hooks ------------------- //

onMounted(() => {
  floorPlanSVG.value = document.getElementById(props.floor);
  floorPlanSVG.value.addEventListener('click', handleClick);
  traverseAndCollectIds(floorPlanSVG.value);
});

onUnmounted(() => {
  floorPlanSVG.value.removeEventListener('click', handleClick);
});
</script>

<style lang="scss" scoped>
.v-menu__content {
  box-shadow: none;
}

.floor-plan__container {
  position: relative;
  /* fill the container, minus the top bar and sc status bar */
  height: calc(100vh - 320px);
  overflow: hidden;
}

.floor-plan__container .pinch-zoom {
  /* fill the container so that zoom controls show in the bottom-right */
  height: 100%;
}

/**
 * This is a custom transition for the menu card above,
 * because the default one is not working properly - on close flies to the top-left corner
*/
.fade-transition-v2 {
  &-leave-active {
    opacity: 0;
  }

  &-enter-active {
    transition: opacity 0.2s ease-in-out;
  }

  &-leave,
  &-leave-to {
    transition: opacity 0s ease-in-out;
  }

  &-enter,
  &-leave-to {
    opacity: 0;
  }
}
</style>
