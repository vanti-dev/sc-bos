<template>
  <content-card>
    <v-data-table
        :headers="headers"
        :items="tableData"
        @click:row="showDevice">
      <template #top/>
    </v-data-table>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import {useDevicesStore} from '@/routes/devices/store';
import {computed, onUnmounted, reactive, ref, watch} from 'vue';
import {usePageStore} from '@/stores/page';

const devicesStore = useDevicesStore();
const pageStore = usePageStore();

const props = defineProps({
  subsystem: {
    type: String,
    default: ''
  }
});

/** @type {Device.Query.AsObject} */
const query = reactive({
  conditionsList: []
});

const headers = ref([
  {text: 'Device name', value: 'name'},
  {text: 'Floor', value: 'metadata.location.floor'},
  {text: 'Zone', value: 'metadata.location.title'}
]);

/** @type {Collection} */
const collection = devicesStore.newCollection();
collection.needsMorePages = true; // todo: this causes us to load all pages, connect with paging logic instead

// watch for changes to the query object and fetch new device list
watch(query, () => collection.query(query), {deep: true, immediate: true});

// watch for changes in the subsystem prop and update query
watch(() => props.subsystem, (sys) => {
  query.conditionsList = query.conditionsList.filter(cond => (cond.field !== 'metadata.membership.subsystem'));
  if (sys !== 'all') {
    query.conditionsList.push({field: 'metadata.membership.subsystem', stringEqual: sys});
  }
}, {immediate: true});

onUnmounted(() => {
  collection.reset(); // stop listening when the component is unmounted
});

const tableData = computed(() => {
  return Object.values(collection.resources.value)
      .map(device => {
        if (device.metadata.location !== undefined && device.metadata.location.moreMap.length > 0) {
          for (const more of device.metadata.location.moreMap) {
            device.metadata.location[more[0]] = more[1];
          }
        }
        return device;
      });
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
