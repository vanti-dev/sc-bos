<template>
  <v-container fluid class="pa-0">
    <v-main class="mx-10 my-6">
      <router-view/>
    </v-main>
    <v-navigation-drawer
        v-if="hasSidebar"
        v-model="showSidebar"
        right
        app
        clipped
        floating
        width="275"
        class="sidebar pa-0">
      <router-view name="sidebar"/>
    </v-navigation-drawer>
  </v-container>
</template>

<script setup>
import {usePage} from '@/components/page';
import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';

const {hasSidebar} = usePage();

const pageStore = usePageStore();
const {showSidebar} = storeToRefs(pageStore);

</script>

<style scoped>
.sidebar{
  background: var(--v-neutral-base);
}
</style>
