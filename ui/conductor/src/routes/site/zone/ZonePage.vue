<template>
  <v-container fluid>
    <v-toolbar flat dense color="transparent">
      <v-btn-toggle mandatory v-model="viewType">
        <v-btn value="map">Map View</v-btn>
        <v-btn value="list">List View</v-btn>
      </v-btn-toggle>
      <v-spacer/>
      <v-btn><v-icon left>mdi-pencil</v-icon>Edit</v-btn>
    </v-toolbar>
    <ZoneMap v-if="viewType === 'map'" :zone="zone"/>
    <ZoneList v-if="viewType === 'list'" :zone="zone"/>
  </v-container>
</template>

<script setup>
import {ref} from 'vue';
import ZoneMap from '@/routes/site/zone/ZoneMap.vue';
import ZoneList from '@/routes/site/zone/ZoneList.vue';

defineProps({
  zone: {
    type: String,
    default: ''
  }
});

const viewType = ref('list');

</script>

<style scoped>

</style>
