<template>
  <content-card>
    <v-data-table
        v-model="selectedDevicesComp"
        :class="tableClasses"
        fixed-header
        :headers="headers"
        hide-default-footer
        :items="tableData"
        :item-class="rowClass"
        item-key="name"
        :items-per-page="devicesPerPage"
        :page.sync="activePage">
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
                  :disabled="floorList.length <= 1"
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
      <template #body="{items}">
        <tbody>
          <tr
              v-for="item in items"
              v-intersect="{
                handler: (entries, observer) => onRowIntersect(entries, observer, item.name),
                options: {
                  rootMargin: '-50px 0px 0px 0px',
                  threshold: 0.75
                }
              }"
              :key="item.name + devicesPerPage"
              @click="showDevice(item)">
            <td>{{ item.metadata.appearance ? item.metadata.appearance.title : item.name }}</td>
            <td>{{ item.metadata?.location?.floor ?? '' }}</td>
            <td>{{ item.metadata?.location?.title ?? '' }}</td>
            <td>
              <WithOccupancy
                  v-if="findSensor(item, 'Occupancy')"
                  class="text-center"
                  :name="item.name"
                  :paused="!intersectedItemNames[item.name]"
                  v-slot="{ occupancyState, occupancyValue }">
                <p :class="[occupancyState.toLowerCase(), 'ma-0 text-body-2']">{{ occupancyState }}</p>
                <v-progress-linear color="primary" indeterminate :active="occupancyValue.loading"/>
              </WithOccupancy>
            </td>
          </tr>
        </tbody>
      </template>
      <template #footer>
        <v-divider/>
        <v-row class="ma-0 pa-0 pt-2 mb-n2">
          <v-spacer/>
          <v-col cols="auto">
            <v-pagination
                v-show="devicesPerPage !== -1"
                :length="pageCount"
                show-current-page
                total-visible="5"
                :value="activePage"
                @input="activePage = $event"/>
          </v-col>
          <v-col cols="auto">
            <v-text-field
                v-show="devicesPerPage !== -1"
                v-model="activePage"
                label="Go to page"
                type="number"
                min="1"
                :max="pageCount"
                outlined
                hide-details
                dense
                style="width: 100px"
                @input="activePage = parseInt($event, 10)"/>
          </v-col>
          <v-col cols="auto">
            <v-select
                dense
                outlined
                hide-details
                :value="devicesPerPage"
                label="Devices per page"
                :items="perPageChoices"
                style="width: 150px; cursor: pointer;"
                @change="devicesPerPage = parseInt($event, 10);"/>
          </v-col>
        </v-row>
      </template>
    </v-data-table>
  </content-card>
</template>

<script setup>
import {computed, onBeforeMount, onMounted, onUnmounted, reactive, ref, watch} from 'vue';
import {storeToRefs} from 'pinia';
import ContentCard from '@/components/ContentCard.vue';
import WithOccupancy from '@/routes/devices/components/renderless-components/WithOccupancy.vue';

import {useErrorStore} from '@/components/ui-error/error';
import {useDevicesStore} from '@/routes/devices/store';
import {Zone} from '@/routes/site/zone/zone';
import {usePageStore} from '@/stores/page';
import {useOccupancyStore} from '@/routes/devices/components/renderless-components/occupancyStore';

const devicesStore = useDevicesStore();
const pageStore = usePageStore();

const {
  findSensor,
  intersectedItemNames,
  onRowIntersect,
  resetIntersectedItemNames
} = useOccupancyStore();
const {activePage, devicesPerPage, perPageChoices} = storeToRefs(useOccupancyStore());

const {floorListResource, handleFloorListLoad} = devicesStore;
const {filterFloor} = storeToRefs(devicesStore);

const errorStore = useErrorStore();

const props = defineProps({
  subsystem: {
    type: String,
    default: 'all'
  },
  zone: {
    type: Zone,
    default: () => {
    }
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
  {text: 'Location', value: 'metadata.location.title'},
  {text: '', value: 'hotpoints', align: 'end', width: '100'}
]);


const search = ref('');


onBeforeMount(() => {
  handleFloorListLoad('pull');
});
onUnmounted(() => {
  handleFloorListLoad('close');
});


const NO_FLOOR = '< no floor >';
const floorList = computed(() => {
  const fieldCounts = floorListResource.value?.fieldCountsList || [];
  const floorFieldCounts = fieldCounts.find(v => v.field === 'metadata.location.floor');
  if (!floorFieldCounts) return [];
  if (floorFieldCounts.countsMap.size <= 0) return [];
  const dst = floorFieldCounts.countsMap.map(([k]) => {
    if (k === '') return NO_FLOOR;
    return k;
  });
  dst.unshift('All');
  return dst;
});

// todo: get this from somewhere. Probably also filter by floor
/* const zoneList = ref([
  'All',
  'L03_013/Meeting Room 1',
  'L03_014/Reception',
  'L03_015/Meeting Room 2'
]);
const filterZone = ref(zoneList.value[0]); */

/** @type {Collection} */
const collection = reactive(devicesStore.newCollection());
collection.needsMorePages = true; // todo: this causes us to load all pages, connect with paging logic instead

// Computed
// ////
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
    emit('update:selectedDevices', value);
  }
});

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
  switch (filterFloor.value.toLowerCase()) {
    case 'all':
      // no filter
      break;
    case NO_FLOOR:
      q.conditionsList.push({field: 'metadata.location.floor', stringEqualFold: ''});
      break;
    default:
      q.conditionsList.push({field: 'metadata.location.floor', stringEqualFold: filterFloor.value});
      break;
  }
  /*   if (filterZone.value.toLowerCase() !== 'all') {
    q.conditionsList.push({field: 'metadata.location.title', stringEqualFold: filterZone.value});
  } */
  return q;
});

const tableData = computed(() => {
  return Object.values(collection.resources.value)
      .filter(props.filter);
});

//
// Table footer
const totalRecords = computed(() => tableData.value.length);
const pageCount = computed(() => Math.ceil(totalRecords.value / devicesPerPage.value));


// Methods
// /////
/**
 * Shows the device in the sidebar
 *
 * @param {*} item
 */
function showDevice(item) {
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


// Watchers
// ////
// watch for changes to the query object and fetch new device list
watch(query, () => collection.query(query.value), {deep: true, immediate: true});
watch(search, () => resetIntersectedItemNames()); // remove old streams


// Lifecycle
// ////
// UI error handling
let unwatchErrors;
onMounted(() => {
  unwatchErrors = errorStore.registerCollection(collection);
});
onUnmounted(() => {
  if (unwatchErrors) unwatchErrors();
  collection.reset(); // stop listening when the component is unmounted
});


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


.v-data-table:not(.selectable) :deep(.v-data-table__selected) {
  background: none;
}

.v-data-table.rowSelectable :deep(.item-selected) {
  background-color: var(--v-primary-darken4);
}

:deep(.v-pagination .v-pagination__item) {
  background-color: var(--v-neutral-lighten1);
}
</style>
