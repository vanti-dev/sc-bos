<template>
  <content-card>
    <v-data-table
        v-model="selectedDevicesComp"
        :headers="headers"
        :items="devicesData"
        item-key="name"
        :item-class="rowClass"
        :footer-props="{
          'items-per-page-options': [
            20,
            50,
            100
          ]
        }"
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
    </v-data-table>
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

const sidebar = useSidebarStore();
const {
  floorList,
  filterFloor,
  search,
  devicesData
} = useDevices(props); // composables/useDevices

const emit = defineEmits(['update:selectedDevices']);

const headers = ref([
  {value: 'metadata.membership.subsystem', width: '20px', class: 'pl-4 pr-0', cellClass: 'pl-4 pr-0', sortable: false},
  {text: 'Device name', value: 'name'},
  {text: 'Floor', value: 'metadata.location.floor'},
  {text: 'Zone', value: 'metadata.location.zone'},
  {text: '', value: 'hotpoint', align: 'end', width: '100', sortable: false}
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
 * @param {*} item
 */
function showDevice(item) {
  sidebar.visible = true;
  sidebar.title = item.metadata.appearance ? item.metadata.appearance.title : item.name;
  sidebar.data = item;
}

/**
 * @param {*} item
 * @return {string}
 */
function rowClass(item) {
  if (sidebar.visible && sidebar.data?.name === item.name) {
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
