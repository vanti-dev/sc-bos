<template>
  <content-card>
    <v-data-table
        :headers="headers"
        :items="tableData"
        @click:row="showDevice">
      <template #top/>
      <template #item.setPoint="{item}">{{ getSetPoint(item.name) }}</template>
    </v-data-table>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import {useHvacStore} from '@/routes/devices/hvac/store';
import {storeToRefs} from 'pinia';
import {computed, onUnmounted, reactive, ref, watch} from 'vue';
import {usePageStore} from '@/stores/page';
import {timestampToDate} from '@/api/convpb';

const hvacStore = useHvacStore();
const {getSetPoint} = storeToRefs(hvacStore);

const pageStore = usePageStore();

/** @type {Device.Query.AsObject} */
const query = reactive({
  conditionsList: [
    {field: 'metadata.membership.subsystem', stringEqual: 'HVAC'}
  ]
});

const headers = ref([
  {text: 'Device name', value: 'name'},
  {text: 'Floor', value: 'metadata.location.floor'},
  {text: 'Zone', value: 'metadata.location.title'},
  {text: 'Device GUID', value: 'guid'}
]);

/** @type {Collection} */
const collection = hvacStore.newCollection();
collection.needsMorePages = true; // todo: this causes us to load all pages, connect with paging logic instead

watch(query, () => collection.query(query), {deep: true, immediate: true});

onUnmounted(() => {
  collection.reset(); // stop listening when the component is unmounted
});

const tableData = computed(() => {
  const devices = Object.values(collection.resources.value)
      .map(device => {
        if (device.metadata.location !== undefined && device.metadata.location.moreMap.length > 0) {
          for (const more of device.metadata.location.moreMap) {
            device.metadata.location[more[0]] = more[1];
          }
        }
        return device;
      });
  return devices;
});

/**
 * Shows the device in the sidebar
 *
 * @param {*} item
 * @param {*} row
 */
function showDevice(item, row) {
  pageStore.showSidebar = true;
  pageStore.sidebarTitle = item.name;
  pageStore.sidebarData = item;
}

</script>

<style scoped>

</style>
