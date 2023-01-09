<template>
  <content-card>
    <v-data-table
        :headers="headers"
        :items="deviceList"
        @click:row="showDevice">
      <template #top/>
      <template #item.setPoint="{item}">{{ getSetPoint(item.deviceId) }}</template>
    </v-data-table>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import {useHvacStore} from '@/routes/devices/hvac/store';
import {storeToRefs} from 'pinia';
import {ref} from 'vue';
import {usePageStore} from '@/stores/page';

const hvacStore = useHvacStore();
const {deviceList, getSetPoint} = storeToRefs(hvacStore);

const pageStore = usePageStore();

const headers = ref([
  {text: 'Device name', value: 'deviceId'},
  {text: 'Floor', value: 'floor'},
  {text: 'Zone(s)', value: 'locations'},
  {text: 'Set Point', value: 'setPoint'},
  {text: 'Device GUID', value: 'guid'}
]);

/**
 * Shows the device in the sidebar
 *
 * @param {*} item
 * @param {*} row
 */
function showDevice(item, row) {
  console.log(item, row);
  pageStore.showSidebar = true;
  pageStore.sidebarData = item;
}

</script>

<style scoped>

</style>
