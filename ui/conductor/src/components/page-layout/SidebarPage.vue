<template>
  <div>
    <v-main>
      <v-container fluid class="pa-6">
        <router-view/>
      </v-container>
    </v-main>
    <v-navigation-drawer
        v-if="hasSidebar"
        v-model="sidebar.visible"
        ref="sidebarDOMElement"
        class="sidebarDOMElement pa-0"
        :class="{resizing}"
        color="neutral"
        floating
        location="right"
        :width="sideBarWidth">
      <router-view name="sidebar"/>
    </v-navigation-drawer>
  </div>
</template>

<script setup>
import {usePage} from '@/components/page';
import {useSidebarStore} from '@/stores/sidebar';
import {onUnmounted, ref, watchEffect} from 'vue';

const props = defineProps({
  minWidth: {
    type: Number,
    default: 275
  },
  maxWidth: {
    type: Number,
    default: 600
  }
});

const {hasSidebar} = usePage();
const sidebarDOMElement = ref(null);
const drawerBorder = ref(/** @type {HTMLDivElement} */ null);
const sideBarWidth = ref(props.minWidth);

const resizing = ref(false);
const handleShift = ref(0);
const pointerId = ref(0);

const sidebar = useSidebarStore();

const beginDrag = (e) => {
  resizing.value = true; // update styles
  pointerId.value = e.pointerId; // for cleaning up later if needed
  // receive mouse events even when mouse is outside the element/browser window
  drawerBorder.value.setPointerCapture(e.pointerId);
  // record where in the border the mouse was pressed to anchor it there to avoid jumps
  handleShift.value = e.clientX - drawerBorder.value.getBoundingClientRect().left;
  drawerBorder.value.addEventListener('pointermove', drag, false); // Add event listener
};

const drag = (event) => {
  // the width of the sidebar is how far from the right edge of the viewport the mouse is
  let clamped = document.body.scrollWidth - event.clientX + handleShift.value;
  clamped = Math.min(clamped, props.maxWidth);
  clamped = Math.max(clamped, props.minWidth);
  sideBarWidth.value = clamped;
};

const endDrag = (e) => {
  resizing.value = false;
  drawerBorder.value.releasePointerCapture(e.pointerId);
  drawerBorder.value.removeEventListener('pointermove', drag, false);
};

// Set event listeners for sidebar resizing
const setUp = () => {
  if (hasSidebar.value && sidebarDOMElement.value.$el) {
    sideBarWidth.value = props.minWidth; // Set sidebar width to default

    drawerBorder.value = sidebarDOMElement.value.$el.querySelector(
        '.v-navigation-drawer__border'
    ); // Get border element

    drawerBorder.value.addEventListener('pointerdown', beginDrag);
    drawerBorder.value.addEventListener('pointerup', endDrag);
  }
};

const cleanUp = () => {
  // clean up listeners and event state
  drawerBorder.value?.removeEventListener('pointerdown', beginDrag, false);
  drawerBorder.value?.removeEventListener('pointerup', endDrag, false);
  drawerBorder.value?.removeEventListener('pointermove', drag, false);
  drawerBorder.value?.releasePointerCapture(pointerId.value);

  resizing.value = false; // Reset resizing styles
  pointerId.value = 0; // Reset pointer ID
  sidebar.closeSidebar(); // Close sidebar
  sideBarWidth.value = props.minWidth; // Reset sidebar width
};

// Watch for sidebar DOM element
watchEffect(() => {
  if (sidebarDOMElement.value?.$el) {
    setUp(); // Set event listeners
  }
});
onUnmounted(() => {
  cleanUp(); // Remove event listeners
});
</script>

<style lang="scss">
.sidebarDOMElement > .v-navigation-drawer__border {
  width: 4px;
  background-color: rgb(var(--v-theme-primary-darken-4));
  transition: all 0.2s ease-in-out;
  cursor: ew-resize;

  &:hover {
    width: 8px;
    background-color: rgb(var(--v-theme-primary-darken-1));
  }
}

.sidebarDOMElement.resizing {
  transition: none;

  > .v-navigation-drawer__border {
    background-color: rgb(var(--v-theme-primary-darken-1));
    width: 8px;
  }
}
</style>
