<template>
  <content-card>
    <v-data-table
        v-model="selectedDevicesComp"
        :headers="headers"
        :items="tableData"
        item-key="name"
        :item-class="rowClass"
        :show-select="showSelect"
        :class="tableClasses"
        @click:row="showDevice">
      <template #top>
        <!-- todo: bulk actions -->
        <!-- filters -->
        <v-container fluid style="width: 100%">
          <v-row dense>
            <v-col cols="12" md="5">
              <v-text-field
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
            <!--            <v-col cols="12" md="2">
              <v-select
                  v-model="filterZone"
                  :items="zoneList"
                  label="Zone"
                  hide-details
                  filled/>
            </v-col>-->
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
import {Zone} from '@/routes/site/zone/zone';

const devicesStore = useDevicesStore();
const pageStore = usePageStore();

const props = defineProps({
  subsystem: {
    type: String,
    default: 'all'
  },
  zone: {
    type: Zone,
    default: () => {}
  },
  showSelect: {
    type: Boolean,
    default: true
  },
  rowSelect: {
    type: Boolean,
    default: true
  },
  selectedDevices: {
    type: Array,
    default: () => []
  },
  filter: {
    type: Function,
    default: () => true
  }
});

const emit = defineEmits(['update:selectedDevices']);

const headers = ref([
  {text: 'Device name', value: 'name'},
  {text: 'Floor', value: 'metadata.location.floor'},
  {text: 'Location', value: 'metadata.location.title'}
]);

const tableClasses = computed(() => {
  const c = [];
  if (props.showSelect) c.push('selectable');
  if (props.rowSelect) c.push('rowSelectable');
  return c.join(' ');
});

const selectedDevicesComp = computed({
  get() {
    return tableData.value.filter(device => props.selectedDevices.indexOf(device.name) >= 0);
  },
  set(value) {
    emit('update:selectedDevices', value.map(d => d.name));
  }
});

const search = ref('');

// todo: get this from somewhere
const floorList = ref([
  'All', 'L00', 'L01', 'L02', 'L03', 'L04'
]);
const filterFloor = ref(floorList.value[0]);

// todo: get this from somewhere. Probably also filter by floor
/* const zoneList = ref([
  'All',
  'L03_013/Meeting Room 1',
  'L03_014/Reception',
  'L03_015/Meeting Room 2'
]);
const filterZone = ref(zoneList.value[0]); */

/** @type {Collection} */
const collection = devicesStore.newCollection();
collection.needsMorePages = true; // todo: this causes us to load all pages, connect with paging logic instead

/** @type {ComputedRef<Device.Query.AsObject>} */
const query = computed(() => {
  const q = {conditionsList: []};
  if (search.value) {
    const words = search.value.split(/\s+/);
    q.conditionsList.push(...words.map(word => ({stringContainsFold: word})));
  }
  if (props.subsystem.toLowerCase() !== 'all') {
    q.conditionsList.push({field: 'metadata.membership.subsystem', stringEqualFold: props.subsystem});
  }
  /*   if (filterZone.value.toLowerCase() !== 'all') {
    q.conditionsList.push({field: 'metadata.location.title', stringEqualFold: filterZone.value});
  } */
  return q;
});

// watch for changes to the query object and fetch new device list
watch(query, () => collection.query(query.value), {deep: true, immediate: true});

onUnmounted(() => {
  collection.reset(); // stop listening when the component is unmounted
});

const tableData = computed(() => {
  return Object.values(collection.resources.value)
      .filter(props.filter)
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

/**
 * @param {*} item
 * @return {string}
 */
function rowClass(item) {
  if (pageStore.showSidebar && pageStore.sidebarData?.name === item.name) {
    return 'item-selected';
  }
  return '';
}

</script>

<style lang="scss" scoped>
:deep(.v-data-table-header__icon) {
  margin-left: 8px;
}

.v-data-table :deep(.v-data-footer) {
  background: var(--v-neutral-lighten1) !important;
  border-radius: 0px 0px $border-radius-root*2 $border-radius-root*2;
  border: none;
  margin: 0 -12px -12px;
}


.v-data-table:not(.selectable) :deep(.v-data-table__selected){
  background: none;
}

.v-data-table.rowSelectable :deep(.item-selected) {
  background-color: var(--v-primary-darken4);
}
</style>
