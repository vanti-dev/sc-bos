<template>
  <v-container fluid class="pa-0">
    <v-main class="mx-10 my-6">
      <router-view/>
    </v-main>
    <v-navigation-drawer
        v-if="hasSidebar"
        v-model="showSidebar"
        id="right-sidebar"
        app
        class="pa-0"
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
import {ref, onMounted} from 'vue';
import {storeToRefs} from 'pinia';
import {usePage} from '@/components/page';
import {usePageStore} from '@/stores/page';

const {hasSidebar} = usePage();

const pageStore = usePageStore();
const {showSidebar} = storeToRefs(pageStore);

const sideBar = ref({
  width: 275,
  borderSize: 4
});

/**
 * Set the initial border width and styles of the sidebar.
 */
const setBorderWidth = () => {
  const mainElement = document.querySelector('#right-sidebar');
  const innerElement = mainElement.querySelector('.v-navigation-drawer__border');
  innerElement.style.width = sideBar.value.borderSize + 'px';
  innerElement.style.cursor = 'ew-resize';
  innerElement.style.backgroundColor = 'var(--v-primary-darken4)';
};

/**
 * Set event listeners for sidebar resizing.
 */
const setEvents = () => {
  const minSize = sideBar.value.borderSize;
  const mainElement = document.querySelector('#right-sidebar');
  const drawerBorder = mainElement.querySelector('.v-navigation-drawer__border');
  const direction = mainElement.classList.contains('v-navigation-drawer--right') ?
    'right' :
    'left';

  /**
   * Handle the resizing of the sidebar.
   *
   * @param {MouseEvent} event - The mouse event object.
   */
  const resize = (event) => {
    document.body.style.cursor = 'ew-resize';
    let forced =
      direction === 'right' ?
        document.body.scrollWidth - event.clientX :
        event.clientX;

    // Enforce maximum width of 600px
    forced = Math.min(forced, 600);

    // Enforce minimum width of 275px
    forced = Math.max(forced, 275);

    mainElement.style.width = forced + 'px';
  };

  drawerBorder.addEventListener('mousedown', (event) => {
    if (event.offsetX < minSize) {
      mainElement.style.transition = 'initial';
      document.addEventListener('mousemove', resize, false);
    }
  });

  document.addEventListener('mouseup', () => {
    mainElement.style.transition = '';
    sideBar.value.width = mainElement.style.width;
    document.body.style.cursor = '';
    document.removeEventListener('mousemove', resize, false);
  });
};

onMounted(() => {
  if (hasSidebar.value) {
    setBorderWidth();
    setEvents();
  }
});
</script>
