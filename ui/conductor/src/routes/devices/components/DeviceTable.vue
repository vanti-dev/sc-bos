<template>
  <content-card>
    <v-data-table-server
        v-model="selectedDevicesComp"
        :headers="headers"
        :items="pagedItems"
        item-key="name"
        :row-props="rowProps"
        @update:options="fetchMoreItems"
        :items-length="totalItems"
        v-model:page="currentPage"
        v-model:items-per-page="itemsPerPage"
        :loading="loading"
        :items-per-page-options="[
          {title: '20', value: 20},
          {title: '50', value: 50},
          {title: '100', value: 100}
        ]"
        :show-select="showSelect"
        item-value="name"
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
                  append-inner-icon="mdi-magnify"
                  label="Search devices"
                  hide-details
                  variant="filled"/>
            </v-col>
            <v-spacer/>
            <v-col cols="12" md="2">
              <v-select
                  :disabled="floorList.length <= 1"
                  v-model="selectedFloor"
                  :items="floorList"
                  label="Floor"
                  hide-details
                  variant="filled"/>
            </v-col>
            <!--            <v-col cols="12" md="2">
              <v-select
                  v-model="filterZone"
                  :items="zoneList"
                  label="Zone"
                  hide-details
                  variant="filled"/>
            </v-col>-->
          </v-row>
        </v-container>
      </template>
      <template #item.metadata.membership.subsystem="{ item }">
        <subsystem-icon size="20px" :subsystem="item.metadata?.membership?.subsystem" no-default/>
      </template>
      <template #item.name="{ item }">
        {{ item.metadata.appearance ? item.metadata.appearance.title : item.name }}
      </template>
      <template #item.hotpoint="{item}">
        <hot-point
            v-slot="{live}"
            class="d-flex align-center justify-end"
            :item-key="item.name"
            style="height:100%">
          <device-cell :paused="!live" :item="item"/>
        </hot-point>
      </template>
    </v-data-table-server>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import HotPoint from '@/components/HotPoint.vue';
import SubsystemIcon from '@/components/SubsystemIcon.vue';
import useDevices from '@/composables/useDevices';
import {Zone} from '@/routes/site/zone/zone';
import {useSidebarStore} from '@/stores/sidebar';
import {computed, ref} from 'vue';
import DeviceCell from './DeviceCell.vue';

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

const sidebar = useSidebarStore();
const search = ref('');
const selectedFloor = ref('All');
const wantCount = ref(20); // same as initial itemsPerPage
const useDevicesOpts = computed(() => {
  return {
    filter: props.filter,
    subsystem: props.subsystem,
    search: search.value,
    floor: selectedFloor.value,
    wantCount: wantCount.value
  };
});
const {
  floorList,
  devicesData,
  totalItems,
  loading
} = useDevices(useDevicesOpts); // composables/useDevices

const currentPage = ref(1);
const itemsPerPage = ref(20);
const fetchMoreItems = ({page, itemsPerPage}) => {
  wantCount.value = page * itemsPerPage;
};
const pagedItems = computed(() => {
  return devicesData.value.slice((currentPage.value - 1) * itemsPerPage.value, currentPage.value * itemsPerPage.value);
});

const headers = ref([
  {key: 'metadata.membership.subsystem', width: '20px', class: 'pl-4 pr-0', cellClass: 'pl-4 pr-0', sortable: false},
  {title: 'Device name', key: 'name'},
  {title: 'Floor', key: 'metadata.location.floor'},
  {title: 'Zone', key: 'metadata.location.zone'},
  {key: 'hotpoint', align: 'end', width: '100', sortable: false}
]);

const tableClasses = computed(() => {
  const c = [];
  if (props.showSelect) c.push('selectable');
  if (props.rowSelect) c.push('rowSelectable');
  return c.join(' ');
});

const selectedDevicesComp = computed({
  get() {
    return devicesData.value.filter(device => props.selectedDevices.indexOf(device.name) >= 0);
  },
  set(value) {
    emit('update:selectedDevices', value);
  }
});

/**
 * Shows the device in the sidebar
 *
 * @param {PointerEvent} e
 * @param {*} item
 */
function showDevice(e, {item}) {
  sidebar.visible = true;
  sidebar.title = item.metadata.appearance ? item.metadata.appearance.title : item.name;
  sidebar.data = item;
}

/**
 * @param {*} item
 * @return {Record<string, any>}
 */
function rowProps({item}) {
  if (sidebar.visible && sidebar.data?.name === item.name) {
    return {class: 'item-selected'};
  }
  return {};
}

</script>

<style lang="scss" scoped>
@use 'vuetify/settings';

:deep(.v-data-table-header__icon) {
  margin-left: 8px;
}

.v-data-table :deep(.v-data-footer) {
  background: rgb(var(--v-theme-neutral-lighten-1)) !important;
  border-radius: 0 0 settings.$border-radius-root*2 settings.$border-radius-root*2;
  border: none;
  margin: 0 -12px -12px;
}


.v-data-table:not(.selectable) :deep(.v-data-table__selected) {
  background: none;
}

.v-data-table.rowSelectable :deep(.item-selected) {
  background-color: rgb(var(--v-theme-primary-darken-4));
}
</style>
