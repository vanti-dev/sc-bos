<template>
  <side-bar>
    <DeviceSideBarContent :device-id="deviceId" :traits="traits"/>
  </side-bar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';
import DeviceSideBarContent from '@/routes/devices/components/DeviceSideBarContent.vue';
import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';
import {computed} from 'vue';

const pageStore = usePageStore();
const {sidebarData} = storeToRefs(pageStore);

const deviceId = computed(() => {
  return sidebarData.value?.name ?? '';
});

const traits = computed(() => {
  const traits = {};
  if (sidebarData.value?.metadata?.traitsList) {
    // flatten array of trait objects (e.g. [{name: 'trait1', ...}, ...] into object (e.g. {trait1: true, ...})
    sidebarData.value.metadata.traitsList.forEach((trait) => traits[trait.name] = true);
  }
  return traits;
});

</script>

<style scoped>
</style>
