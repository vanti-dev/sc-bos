<template>
  <content-card>
    <DataTable
        :dropdown="{
          dropdownItems: floorList,
          dropdownLabel: 'Floor',
          dropdownValue: filterFloor
        }"
        :table-headers="tableHeaders"
        :table-items="tableData"
        :row-select="props.rowSelect"
        :selected-items="computeSelectedDevices"
        :zone="props.zone"
        @update:dropdownValue="filterFloor = $event"
        @update:selectedItems="emit('update:selectedDevices', $event)">
      <template #hotpoint="{findSensor, item, intersectedItemNames}">
        <WithHotpoint
            v-if="findSensor(item, 'Occupancy')"
            class="text-center"
            :name="item.name"
            :paused="!intersectedItemNames[item.name]"
            style="min-width: 75px;">
          <template #occupancy="{ occupancyData }">
            <p :class="[occupancyData.occupancyState.toLowerCase(), 'ma-0 text-body-2']">
              {{ occupancyData.occupancyState }}
            </p>
            <v-progress-linear color="primary" indeterminate :active="occupancyData.occupancyValue.loading"/>
          </template>
        </WithHotpoint>
      </template>
    </DataTable>
  </content-card>
</template>

<script setup>
import {computed, onBeforeMount, onMounted, onUnmounted, watch} from 'vue';
import {storeToRefs} from 'pinia';

// Component import
import ContentCard from '@/components/ContentCard.vue';
import DataTable from '@/components/composables/DataTable/DataTable.vue';
import WithHotpoint from '@/routes/devices/components/renderless-components/WithHotpoint.vue';

// Store imports
import {useErrorStore} from '@/components/ui-error/error';
import {useDevicesStore} from '@/routes/devices/store';
import {useTableDataStore} from '@/stores/tableDataStore';

// Type imports
import {Zone} from '@/routes/site/zone/zone';

const tableDataStore = useTableDataStore();
const {resetIntersectedItemNames, siteEditor} = tableDataStore;
const {search} = storeToRefs(tableDataStore);
const devicesStore = useDevicesStore();
const {collection, floorListResource, handleFloorListLoad} = devicesStore;
const {filterFloor} = storeToRefs(devicesStore);

const errorStore = useErrorStore();

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

onBeforeMount(() => {
  handleFloorListLoad('pull');
});
onUnmounted(() => {
  handleFloorListLoad('close');
});

// Computed
// ////
const tableHeaders = computed(() => {
  const headers = [
    {text: 'Device name', value: 'name'},
    {text: 'Floor', value: 'metadata.location.floor'},
    {text: 'Location', value: 'metadata.location.title'},
    !siteEditor.zone ?? {text: '', value: 'hotpoints', align: 'end', width: '100'}
  ];
  return headers.filter(header => header);
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

collection.needsMorePages = true; // todo: this causes us to load all pages, connect with paging logic instead
const tableData = computed(() => {
  return Object.values(collection.resources.value)
      .filter(props.filter);
});

const computeSelectedDevices = computed(() => {
  return tableData.value.filter(device => props.selectedDevices.indexOf(device.name) >= 0);
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

</style>
