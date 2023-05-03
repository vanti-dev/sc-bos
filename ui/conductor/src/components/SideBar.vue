<template>
  <v-list tile class="ma-0 pa-0">
    <v-subheader class="text-body-large font-weight-bold sidebar-title pa-4 neutral lighten-1">
      <span class="sidebar-title-text my-auto">
        {{ sidebarTitle }}
      </span><v-spacer/>
      <v-btn icon small @click="pageStore.closeSidebar()" color="neutral lighten-4">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-subheader>
    <slot/>
  </v-list>
</template>

<script setup>
import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';

const pageStore = usePageStore();
const {sidebarTitle} = storeToRefs(pageStore);

</script>

<style scoped>
.sidebar-title {
  height: auto;
  align-items: start;
}
.sidebar-title-text {
  display: block;
  overflow-wrap: break-word;
  max-width: 210px;
}
</style>
