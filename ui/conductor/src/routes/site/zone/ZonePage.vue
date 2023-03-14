<template>
  <v-container fluid>
    <v-toolbar flat dense color="transparent" class="mb-3">
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
import {computed, ref} from 'vue';
import ZoneMap from '@/routes/site/zone/ZoneMap.vue';
import {useServicesStore} from '@/stores/services';
import {ServiceNames} from '@/api/ui/services';
import {Zone} from '@/routes/site/zone/zone';
import DeviceTable from '@/routes/devices/components/DeviceTable.vue';
import {Service} from '@sc-bos/ui-gen/proto/services_pb';
import {newActionTracker} from '@/api/resource';

const servicesStore = useServicesStore();
const zoneCollection = ref(servicesStore.getService(ServiceNames.Zones).servicesCollection);

const props = defineProps({
  zone: {
    type: String,
    default: ''
  }
});

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
