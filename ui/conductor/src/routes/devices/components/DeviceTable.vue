<template>
  <content-card>
    <v-data-table
        v-model="selectedDevices"
        :headers="headers"
        :items="tableData"
        item-key="name"
        show-select
        @click:row="showDevice">
      <template #top>
        <!-- todo: bulk actions -->
        <!-- filters -->
        <v-container fluid style="width: 100%">
          <v-row dense>
            <v-col cols="12" md="5">
              <v-text-field
                  disabled
                  v-model="search"
                  append-icon="mdi-magnify"
                  label="Search devices"
                  hide-details
                  filled/>
            </v-col>
            <v-spacer/>
            <v-col cols="12" md="2">
              <v-select
                  disabled
                  v-model="filterFloor"
                  :items="floorList"
                  label="Floor"
                  hide-details
                  filled/>
            </v-col>
            <v-col cols="12" md="2">
              <v-select
                  v-model="filterZone"
                  :items="zoneList"
                  label="Zone"
                  hide-details
                  filled/>
            </v-col>
          </v-row>
        </v-container>
      </template>
    </v-data-table>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import {useDevicesStore} from '@/routes/devices/store';
import {computed, onUnmounted, ref, watch} from 'vue';
import {usePageStore} from '@/stores/page';

const devicesStore = useDevicesStore();
const pageStore = usePageStore();

const props = defineProps({
  subsystem: {
    type: String,
    default: ''
  }
});

const headers = ref([
  {text: 'Device name', value: 'name'},
  {text: 'Floor', value: 'metadata.location.floor'},
  {text: 'Zone', value: 'metadata.location.title'}
]);

const selectedDevices = ref([]);
const search = ref('');

// todo: get this from somewhere
const floorList = ref([
  'All', 'L00', 'L01', 'L02', 'L03', 'L04'
]);
const filterFloor = ref(floorList.value[0]);

// todo: get this from somewhere. Probably also filter by floor
const zoneList = ref([
  'All',
  'L03_013/Meeting Room 1',
  'L03_014/Reception',
  'L03_015/Meeting Room 2'
]);
const filterZone = ref(zoneList.value[0]);

/** @type {Collection} */
const collection = devicesStore.newCollection();
collection.needsMorePages = true; // todo: this causes us to load all pages, connect with paging logic instead

/** @type {ComputedRef<Device.Query.AsObject>} */
const query = computed(() => {
  const q = {conditionsList: []};
  if (props.subsystem.toLowerCase() !== 'all') {
    q.conditionsList.push({field: 'metadata.membership.subsystem', stringEqual: props.subsystem});
  }
  if (filterZone.value.toLowerCase() !== 'all') {
    q.conditionsList.push({field: 'metadata.location.title', stringEqual: filterZone.value});
  }
  return q;
});

// watch for changes to the query object and fetch new device list
watch(query, () => collection.query(query.value), {deep: true, immediate: true});

onUnmounted(() => {
  collection.reset(); // stop listening when the component is unmounted
});

const tableData = computed(() => {
  return Object.values(collection.resources.value)
      .map(device => {
        if (device.metadata?.location?.moreMap?.length > 0) {
          // flatten the moreMap of location to expose additional location data to the user - this allows us to add
          // extra info that isn't natively supported by the smart core metadata trait
          for (const [key, val] of device.metadata.location.moreMap) {
            device.metadata.location[key] = val;
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

<style lang="scss" scoped>
::v-deep(.v-data-table-header__icon) {
  margin-left: 8px;
}

.v-data-table ::v-deep(.v-data-footer) {
  background: var(--v-neutral-lighten1) !important;
  border-radius: 0px 0px $border-radius-root*2 $border-radius-root*2;
  border: none;
  margin: 0 -12px -12px;
}
</style>
