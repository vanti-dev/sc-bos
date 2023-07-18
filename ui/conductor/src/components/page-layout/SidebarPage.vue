<template>
  <v-container fluid class="pa-0">
    <v-main class="mx-10 my-6">
      <router-view/>
    </v-main>
    <v-navigation-drawer
        v-if="hasSidebar"
        v-model="showSidebar"
        ref="rightSidebar"
        app
        class="rightSidebar pa-0"
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
import {ref, onMounted, onBeforeUnmount} from 'vue';
import {storeToRefs} from 'pinia';
import {usePage} from '@/components/page';
import {usePageStore} from '@/stores/page';

const {hasSidebar} = usePage();

const pageStore = usePageStore();
const {showSidebar} = storeToRefs(pageStore);

const rightSidebar = ref(null);
const drawerBorder = ref(null);
const sideBar = ref({
  width: 275,
  borderSize: 4
});

let originalCursor = '';

// Handle mouse move for sidebar resizing
const handleMouseMove = (event) => {
  document.body.style.cursor = 'ew-resize'; // Set cursor
  drawerBorder.value.style.width = '8px'; // Set border width
  let forced = document.body.scrollWidth - event.clientX; // Calculate width

  forced = Math.min(forced, 600); // Set max width
  forced = Math.max(forced, 275); // Set min width

  rightSidebar.value.$el.style.width = `${forced}px`; // Set width
};

const handleMouseUp = () => {
  rightSidebar.value.$el.style.transition = ''; // Remove transition
  sideBar.value.width = rightSidebar.value.$el.style.width; // Update sidebar width
  drawerBorder.value.style.backgroundColor = ''; // Remove border highlight on button release
  drawerBorder.value.style.width = ''; // Reset border width
  document.removeEventListener('mousemove', handleMouseMove, false); // Remove event listener

  document.body.style.cursor = originalCursor; // Restore original cursor value
  rightSidebar.value.$el.style.userSelect = ''; // Enable text selection
};

// Set event listeners for sidebar resizing
const setEvents = () => {
  drawerBorder.value.addEventListener('mousedown', (event) => {
    originalCursor = originalCursor || document.body.style.cursor; // Store original cursor value
    rightSidebar.value.$el.style.transition = 'all 0.1s ease'; // Add transition
    drawerBorder.value.style.backgroundColor = 'var(--v-primary-darken1)'; // Highlight border while moving
    drawerBorder.value.style.width = '8px'; // Set border width
    document.addEventListener('mousemove', handleMouseMove, false); // Add event listener

    rightSidebar.value.$el.style.userSelect = 'none'; // Disable text selection
  });


  document.addEventListener('mouseup', handleMouseUp); // Add event listener
};

onMounted(() => {
  // If sidebar exists and is mounted
  if (hasSidebar.value && rightSidebar.value.$el) {
    drawerBorder.value = rightSidebar.value.$el.querySelector('.v-navigation-drawer__border'); // Get border element
    setEvents(); // Set event listeners
  }
});

onBeforeUnmount(() => {
  document.removeEventListener('mousemove', handleMouseMove, false); // Remove event listener
  document.removeEventListener('mouseup', handleMouseUp, false); // Remove event listener
});

</script>

<style lang="scss">
.rightSidebar > .v-navigation-drawer__border {
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
