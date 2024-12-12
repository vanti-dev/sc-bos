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
        class="resizable pa-0"
        :class="{resizing}"
        color="neutral"
        floating
        location="right"
        :width="sideBarWidth">
      <router-view name="sidebar"/>
      <div class="resize--handle" ref="resizeHandleElement"/>
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
const resizeHandleElement = ref(/** @type {HTMLDivElement} */ null);
const sideBarWidth = ref(props.minWidth);

const resizing = ref(false);
const handleShift = ref(0);
const pointerId = ref(0);

const sidebar = useSidebarStore();

const beginDrag = (e) => {
  resizing.value = true; // update styles
  pointerId.value = e.pointerId; // for cleaning up later if needed
  // receive mouse events even when mouse is outside the element/browser window
  resizeHandleElement.value.setPointerCapture(e.pointerId);
  // record where in the border the mouse was pressed to anchor it there to avoid jumps
  handleShift.value = e.clientX - resizeHandleElement.value.getBoundingClientRect().left;
  resizeHandleElement.value.addEventListener('pointermove', drag, false); // Add event listener
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
  resizeHandleElement.value.releasePointerCapture(e.pointerId);
  resizeHandleElement.value.removeEventListener('pointermove', drag, false);
};

// Set event listeners for sidebar resizing
const setUp = () => {
  if (hasSidebar.value && resizeHandleElement.value) {
    sideBarWidth.value = props.minWidth; // Set sidebar width to default

    resizeHandleElement.value.addEventListener('pointerdown', beginDrag);
    resizeHandleElement.value.addEventListener('pointerup', endDrag);
  }
};

const cleanUp = () => {
  // clean up listeners and event state
  resizeHandleElement.value?.removeEventListener('pointerdown', beginDrag, false);
  resizeHandleElement.value?.removeEventListener('pointerup', endDrag, false);
  resizeHandleElement.value?.removeEventListener('pointermove', drag, false);
  resizeHandleElement.value?.releasePointerCapture(pointerId.value);

  resizing.value = false; // Reset resizing styles
  pointerId.value = 0; // Reset pointer ID
  sidebar.closeSidebar(); // Close sidebar
  sideBarWidth.value = props.minWidth; // Reset sidebar width
};

// Watch for sidebar DOM element
watchEffect(() => {
  if (resizeHandleElement.value) {
    setUp(); // Set event listeners
  }
});
onUnmounted(() => {
  cleanUp(); // Remove event listeners
});
</script>

<style lang="scss">
.resizable .resize--handle {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  width: 4px;
  background-color: rgb(var(--v-theme-primary-darken-4));
  transition: all 0.2s ease-in-out;
  cursor: ew-resize;

  &:hover {
    width: 8px;
    background-color: rgb(var(--v-theme-primary-darken-1));
  }
}

.resizable.resizing {
  transition: none;

   .resize--handle {
    background-color: rgb(var(--v-theme-primary-darken-1));
    width: 8px;
  }
}
</style>
