<template>
  <v-container fluid style="min-height: calc(100vh - 275px); overflow: hidden">
    <PinchZoom @click="handleClick(event)">
      <template #default="{ scale }">
        <!-- eslint-disable-next-line vue/no-v-html -->
        <!-- <div v-html="level" ref="svgContainer" :style="{'--map-scale': scale,...floorPlanStyles}"/> -->
        <Level0 show-doors id="level0" :style="{ '--map-scale': scale }"/>
      </template>
    </PinchZoom>
    <v-menu v-model="showMenu" absolute :position-x="elementWithMenu.x" :position-y="elementWithMenu.y">
      <template #activator="{ on }">
        <div v-if="elementWithMenu.target" v-bind="on" ref="menuActivator"/>
      </template>
      <AccessPointCard/>
    </v-menu>
  </v-container>
</template>

<script setup>
import {onMounted, ref} from 'vue';
import PinchZoom from '@/routes/ops/security/map/PinchZoom.vue';
import Level0 from '@/clients/ew/Level0.vue';
import AccessPointCard from './AccessPointCard.vue';

const svgElement = ref(null);

const elementWithMenu = ref({
  target: '' | null,
  x: 0,
  y: 0
});

const showMenu = ref(false);

const allIds = ref([]);

const traverseAndCollectIds = (element) => {
  if (element.id !== 'level0' && element.tagName !== 'g' && element.id) {
    allIds.value.push(element.id);
  }

  const children = element.children;
  for (let i = 0; i < children.length; i++) {
    traverseAndCollectIds(children[i]);
  }
};

/**
 *
 * @param {HTMLElement} element
 * @return {HTMLElement}
 */
function findDeepestChild(element) {
  let deepestChild = element;

  while (deepestChild.lastElementChild) {
    deepestChild = deepestChild.lastElementChild;
  }

  return deepestChild;
}

/**
 *
 * @param {MouseEvent} event
 */
function handleClick(event) {
  // Reset menu state
  elementWithMenu.value.target = null;
  elementWithMenu.value.x = 0;
  elementWithMenu.value.y = 0;

  // Find deepest child
  const clickedElement = event.target;
  elementWithMenu.value.target = findDeepestChild(clickedElement);

  // Calculate menu position
  const clickedRect = clickedElement.getBoundingClientRect();

  // calculate position to the left
  elementWithMenu.value.x = clickedRect.left - Math.floor((clickedRect.left / 100) * 1.5);
  // calculate position to the bottom
  elementWithMenu.value.y = clickedRect.top + clickedRect.height * 1.5;

  // Show/Hide menu
  showMenu.value = true;
}

onMounted(() => {
  svgElement.value = document.getElementById('level0');
  svgElement.value.addEventListener('click', handleClick);
  traverseAndCollectIds(svgElement.value);
});
</script>
