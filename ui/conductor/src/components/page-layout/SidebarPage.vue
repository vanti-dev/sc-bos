<template>
  <v-container fluid class="pa-0">
    <v-main class="mx-10 my-6">
      <router-view/>
    </v-main>
    <v-navigation-drawer
        v-if="hasSidebar"
        v-model="showSidebar"
        ref="sidebarDOMElement"
        app
        class="sidebarDOMElement pa-0"
        clipped
        color="neutral"
        floating
        right
        :width="sideBar.width">
      <router-view name="sidebar"/>
    </v-navigation-drawer>
  </v-container>
</template>

<script setup>
import {ref, watchEffect} from 'vue';
import {storeToRefs} from 'pinia';
import {usePage} from '@/components/page';
import {usePageStore} from '@/stores/page';

const {hasSidebar} = usePage();
const sidebarDOMElement = ref(null);
const drawerBorder = ref(null);
const sideBar = ref({
  width: 275,
  borderSize: 4
});

const pageStore = usePageStore();
const {showSidebar} = storeToRefs(pageStore);

let originalCursor = '';

const handleMouseDown = (event) => {
  originalCursor = originalCursor || document.body.style.cursor; // Store original cursor value
  sidebarDOMElement.value.$el.style.transition = 'none'; // Add transition
  drawerBorder.value.style.backgroundColor = 'var(--v-primary-darken1)'; // Highlight border while moving
  drawerBorder.value.style.width = '8px'; // Set border width
  document.addEventListener('mousemove', handleMouseMove, false); // Add event listener

  sidebarDOMElement.value.$el.style.userSelect = 'none'; // Disable text selection
};

// Handle mouse move for sidebar resizing
const handleMouseMove = (event) => {
  document.body.style.cursor = 'ew-resize'; // Set cursor
  drawerBorder.value.style.width = '8px'; // Set border width
  let forced = document.body.scrollWidth - event.clientX; // Calculate width

  forced = Math.min(forced, 600); // Set max width
  forced = Math.max(forced, 275); // Set min width

  sidebarDOMElement.value.$el.style.width = `${forced}px`; // Set width
};

const handleMouseUp = () => {
  if (sidebarDOMElement.value && sidebarDOMElement.value.$el) {
    sidebarDOMElement.value.$el.style.transition = ''; // Remove transition
    sideBar.value.width = sidebarDOMElement.value.$el.style.width; // Update sidebar width
  }

  if (drawerBorder.value) {
    drawerBorder.value.style.backgroundColor = ''; // Remove border highlight on button release
    drawerBorder.value.style.width = ''; // Reset border width
  }

  document.removeEventListener('mousemove', handleMouseMove, false); // Remove event listener

  document.body.style.cursor = originalCursor; // Restore original cursor value

  if (sidebarDOMElement.value && sidebarDOMElement.value.$el) {
    sidebarDOMElement.value.$el.style.userSelect = ''; // Enable text selection
  }
};

// Set event listeners for sidebar resizing
const setEvents = () => {
  drawerBorder.value.addEventListener('mousedown', handleMouseDown);
  document.addEventListener('mouseup', handleMouseUp);
};

const setUp = () => {
  if (hasSidebar.value && sidebarDOMElement.value.$el) {
    sideBar.value.width = 275; // Set sidebar width to default

    drawerBorder.value = sidebarDOMElement.value.$el.querySelector(
        '.v-navigation-drawer__border'
    ); // Get border element
    setEvents(); // Set event listeners
  }
};


const cleanUp = () => {
  drawerBorder.value?.removeEventListener('mousedown', handleMouseDown, false); // Remove event listener
  document.removeEventListener('mousemove', handleMouseMove, false); // Remove event listener
  document.removeEventListener('mouseup', handleMouseUp, false); // Remove event listener
  showSidebar.value = false; // Close sidebar
  sideBar.value.width = 275; // Reset sidebar width
  drawerBorder.value = null; // Remove border element leftover
  sidebarDOMElement.value = null; // Remove sidebar DOM element leftover
};

// Watch for sidebar DOM element
watchEffect(() => {
  if (sidebarDOMElement.value?.$el) {
    setUp(); // Set event listeners
  } else {
    cleanUp(); // Remove event listeners
  }
});
</script>

<style lang="scss">
.sidebarDOMElement > .v-navigation-drawer__border {
  width: 4px;
  background-color: var(--v-primary-darken4);
  transition: all 0.2s ease-in-out;
  &:hover {
    width: 8px;
    cursor: ew-resize;
    background-color: var(--v-primary-darken1);
  }
}
</style>
