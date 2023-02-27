<template>
  <v-container fluid>
    <v-toolbar flat dense color="transparent">
      <v-spacer/>
      <v-btn-toggle mandatory dense v-model="viewType">
        <v-btn value="map" disabled>Map View</v-btn>
        <v-btn value="list">List View</v-btn>
      </v-btn-toggle>
      <v-btn class="ml-6" v-if="editMode" @click="editMode=false"><v-icon left>mdi-content-save</v-icon>Save</v-btn>
      <v-btn class="ml-6" v-else @click="editMode=true"><v-icon left>mdi-pencil</v-icon>Edit</v-btn>
    </v-toolbar>
    <div v-if="!zone"/>
    <zone-map v-else-if="viewType === 'map'" :zone="zoneObj"/>
    <device-table
        v-else-if="viewType === 'list'"
        :zone="zoneObj"
        :show-select="editMode"
        :filter="zoneDevicesFilter"/>
  </v-container>
</template>

<script setup>
import {computed, ref} from 'vue';
import ZoneMap from '@/routes/site/zone/ZoneMap.vue';
import {useServicesStore} from '@/stores/services';
import {ServiceNames} from '@/api/ui/services';
import {Zone} from '@/routes/site/zone/zone';
import DeviceTable from '@/routes/devices/components/DeviceTable.vue';

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

const deviceList = computed(() => {
  return zoneObj.value.devices;
});

/**
 *
 * @param device
 */
function zoneDevicesFilter(device) {
  return zoneObj?.value?.deviceIds?.indexOf(device.name) >= 0 ?? true;
}

</script>

<style scoped>

</style>
