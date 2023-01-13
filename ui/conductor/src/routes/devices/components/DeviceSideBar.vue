<template>
  <side-bar>
    <v-list tile class-name="ma-0 pa-0">
      <device-info-card/>
      <air-temperature-card :name="sidebarData.name"/>
    </v-list>
  </side-bar>
</template>

<script setup>
import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';
import AirTemperatureCard from '@/routes/devices/components/trait-cards/AirTemperatureCard.vue';
import DeviceInfoCard from '@/routes/devices/components/trait-cards/DeviceInfoCard.vue';
import SideBar from '@/components/SideBar.vue';

const pageStore = usePageStore();
const {sidebarData} = storeToRefs(pageStore);

</script>

<style scoped>
.sidebarTitle {
  background: var(--v-neutral-lighten1);
  height: auto;
  font-weight: bold;
}
</style>
