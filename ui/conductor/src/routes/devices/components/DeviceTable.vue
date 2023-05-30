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
      <template #item.hotpoint="{item}">
        <HotPoint :item="item" :item-key="item.name">
          <template #hotpoint="{live, sensorTypes}">
            <WithOccupancy
                :name="item.name"
                :paused="!live"
                :sensor-types="sensorTypes">
              <template #occupancy="{occupancyData}">
                <p :class="[occupancyData.occupancyState.toLowerCase(), 'ma-0 text-body-2']">
                  {{ occupancyData.occupancyState }}
                </p>
                <v-progress-linear color="primary" indeterminate :active="occupancyData.occupancyValue.loading"/>
              </template>
            </WithOccupancy>
          </template>
        </HotPoint>
      </template>
    </v-data-table>
  </content-card>
</template>

<script setup>
import {computed, onMounted, onUnmounted, reactive, ref, watch} from 'vue';
import {closeResource, newResourceValue} from '@/api/resource';
import {pullDevicesMetadata} from '@/api/ui/devices';

import {useErrorStore} from '@/components/ui-error/error';
import {useDevicesStore} from '@/routes/devices/store';
import {usePageStore} from '@/stores/page';

import {Zone} from '@/routes/site/zone/zone';

import ContentCard from '@/components/ContentCard.vue';
import WithOccupancy from './renderless/WithOccupancy.vue';
import HotPoint from '@/components/HotPoint.vue';

const devicesStore = useDevicesStore();
const pageStore = usePageStore();
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
    default: false
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
  {text: '', value: 'hotpoint', align: 'end', width: '100'}
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
    emit('update:selectedDevices', value);
  }
});

const search = ref('');

// todo: this information should come from a store and be reusable between components.
const floorListResource = reactive(newResourceValue());
onMounted(() => {
  const req = {includes: {fieldsList: ['metadata.location.floor']}, updatesOnly: false};
  pullDevicesMetadata(req, floorListResource);
});
onUnmounted(() => {
  closeResource(floorListResource);
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
const filterFloor = ref('All');

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

// watch for changes to the query object and fetch new device list
watch(query, () => collection.query(query.value), {deep: true, immediate: true});

// UI error handling
let unwatchErrors;
onMounted(() => {
  unwatchErrors = errorStore.registerCollection(collection);
});
onUnmounted(() => {
  if (unwatchErrors) unwatchErrors();
  collection.reset(); // stop listening when the component is unmounted
});

const tableData = computed(() => {
  return Object.values(collection.resources.value)
      .filter(props.filter);
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


.v-data-table:not(.selectable) :deep(.v-data-table__selected) {
  background: none;
}

.v-data-table.rowSelectable :deep(.item-selected) {
  background-color: var(--v-primary-darken4);
}
</style>
