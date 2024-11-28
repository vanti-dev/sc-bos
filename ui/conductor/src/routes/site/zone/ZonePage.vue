<template>
  <v-container fluid>
    <v-toolbar flat dense color="transparent" class="mb-3">
      <v-combobox
          v-if="configStore.config?.hub"
          v-model="userConfig.node"
          :items="cohort.cohortNodes"
          label="System Component"
          item-title="name"
          item-value="name"
          hide-details="auto"
          :loading="cohort.loading"
          variant="outlined"/>
      <v-spacer/>
      <v-btn class="ml-6" v-if="editMode" @click="save" color="accent" :disabled="blockActions">
        <v-icon start>mdi-content-save</v-icon>
        Save
      </v-btn>
      <v-btn class="ml-6" v-else @click="editMode=true" :disabled="blockActions">
        <v-icon start>mdi-pencil</v-icon>
        Edit
      </v-btn>
    </v-toolbar>
    <div v-if="!zone"/>
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
import {usePullService} from '@/composables/services.js';
import useAuthSetup from '@/composables/useAuthSetup';
import DeviceTable from '@/routes/devices/components/DeviceTable.vue';
import {Zone} from '@/routes/site/zone/zone';
import {useCohortStore} from '@/stores/cohort.js';
import {useUserConfig} from '@/stores/userConfig.js';
import {useUiConfigStore} from '@/stores/uiConfig.js';
import {Service} from '@vanti-dev/sc-bos-ui-gen/proto/services_pb';
import {computed, ref} from 'vue';

const {blockActions} = useAuthSetup();

const userConfig = useUserConfig();
const configStore = useUiConfigStore();
const cohort = useCohortStore();

const props = defineProps({
  zone: {
    type: String,
    default: ''
  }
});

const {value: zoneRes} = usePullService(computed(() => {
  if (!userConfig.node || !props.zone) return null;
  return {
    name: userConfig.node?.name + '/' + ServiceNames.Zones,
    id: props.zone
  };
}), computed(() => ({
  paused: !userConfig.node
})));
const zoneObj = computed(() => {
  const z = zoneRes.value ?? new Service().toObject();
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
