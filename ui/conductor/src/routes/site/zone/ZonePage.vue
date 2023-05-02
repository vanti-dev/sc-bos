<template>
  <v-container fluid>
    <v-toolbar flat dense color="transparent" class="mb-3">
      <v-combobox
          v-if="configStore.config?.hub"
          v-model="node"
          :items="Object.values(hubStore.nodesList)"
          label="System Component"
          item-text="name"
          item-value="name"
          hide-details="auto"
          :loading="hubStore.nodesListCollection.loading ?? true"
          outlined/>
      <v-spacer/>
      <v-btn-toggle mandatory dense v-model="viewType">
        <v-btn value="map" disabled>Map View</v-btn>
        <v-btn value="list">List View</v-btn>
      </v-btn-toggle>
      <v-btn class="ml-6" v-if="editMode" @click="save" color="accent">
        <v-icon left>mdi-content-save</v-icon>
        Save
      </v-btn>
      <v-btn class="ml-6" v-else @click="editMode=true">
        <v-icon left>mdi-pencil</v-icon>
        Edit
      </v-btn>
    </v-toolbar>
    <div v-if="!zone"/>
    <zone-map v-else-if="viewType === 'map'" :zone="zoneObj"/>
    <device-table
        v-else-if="viewType === 'list'"
        :zone="zoneObj"
        :show-select="editMode"
        :row-select="false"
        :filter="zoneDevicesFilter"
        :selected-devices="deviceList"
        @update:selectedDevices="deviceList = $event"/>
  </v-container>
</template>

<script setup>
import {newActionTracker} from '@/api/resource';
import {ServiceNames} from '@/api/ui/services';
import DeviceTable from '@/routes/devices/components/DeviceTable.vue';
import {Zone} from '@/routes/site/zone/zone';
import ZoneMap from '@/routes/site/zone/ZoneMap.vue';
import {useAppConfigStore} from '@/stores/app-config';
import {useHubStore} from '@/stores/hub';
import {usePageStore} from '@/stores/page';
import {useServicesStore} from '@/stores/services';
import {Service} from '@sc-bos/ui-gen/proto/services_pb';
import {computed, ref, watch} from 'vue';

const servicesStore = useServicesStore();
const pageStore = usePageStore();
const configStore = useAppConfigStore();
const hubStore = useHubStore();
const zoneCollection = ref();

const props = defineProps({
  zone: {
    type: String,
    default: ''
  }
});

const node = computed({
  get() {
    return pageStore.sidebarNode;
  },
  set(val) {
    pageStore.sidebarNode = val;
  }
});

watch(node, async () => {
  zoneCollection.value = servicesStore.getService(
      ServiceNames.Zones,
      await node.value.commsAddress,
      await node.value.commsName).servicesCollection;
}, {immediate: true});

const zoneObj = computed(() => {
  const z = zoneCollection?.value?.resources?.value[props.zone] ?? (new Service()).toObject();
  return new Zone(z);
});

const viewType = ref('list');
const editMode = ref(false);

const deviceList = computed({
  get() {
    return zoneObj.value.deviceIds;
  },
  set(value) {
    zoneObj.value.devices = value;
  }
});


/**
 * @param {Device.AsObject} device
 * @return {boolean}
 */
function zoneDevicesFilter(device) {
  return editMode.value || (zoneObj?.value?.deviceIds?.indexOf(device.name) >= 0 ?? true);
}

const saveTracker = newActionTracker();

/**
 * Save the new device list to the zone
 */
function save() {
  zoneObj.value.save(saveTracker);
  editMode.value = false;
}

</script>

<style scoped>

</style>
