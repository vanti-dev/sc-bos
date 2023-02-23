<template>
  <v-container fluid>
    <v-toolbar flat dense color="transparent">
      <v-toolbar-title>{{ zone }}</v-toolbar-title>
      <v-spacer/>
      <v-btn-toggle mandatory dense v-model="viewType">
        <v-btn value="map">Map View</v-btn>
        <v-btn value="list">List View</v-btn>
      </v-btn-toggle>
      <v-spacer/>
      <v-btn><v-icon left>mdi-pencil</v-icon>Edit</v-btn>
    </v-toolbar>
    <div v-if="!zone"/>
    <ZoneMap v-else-if="viewType === 'map'" :zone="zoneObj"/>
    <ZoneList v-else-if="viewType === 'list'" :zone="zoneObj"/>
  </v-container>
</template>

<script setup>
import {computed, ref} from 'vue';
import ZoneMap from '@/routes/site/zone/ZoneMap.vue';
import ZoneList from '@/routes/site/zone/ZoneList.vue';
import {useServicesStore} from '@/stores/services';
import {ServiceNames} from '@/api/ui/services';
import {Zone} from '@/routes/site/zone/zone';

const servicesStore = useServicesStore();
const zoneCollection = ref(servicesStore.getService(ServiceNames.Zones).servicesCollection);

const props = defineProps({
  zone: {
    type: String,
    default: ''
  }
});

const zoneObj = computed(() => {
  const z = zoneCollection?.value?.resources?.value[props.zone] ?? null;
  if (z) {
    return new Zone(z);
  } else {
    return null;
  }
});

const viewType = ref('list');

</script>

<style scoped>

</style>
