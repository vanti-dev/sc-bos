<template>
  <content-card>
    <v-card-title class="pb-0 text-h3">Zone devices</v-card-title>
    <v-data-table :items="selectedDevices" :headers="headers" v-if="zone"/>
    <v-card-text v-else>
      No devices selected
    </v-card-text>
  </content-card>
</template>

<script setup>
import {Zone} from '@/routes/site/zone/zone';
import ContentCard from '@/components/ContentCard.vue';
import {computed} from 'vue';

const props = defineProps({
  zone: {
    type: Zone,
    default: () => {
      return null;
    }
  },
  editMode: Boolean
});

const headers = [
  {text: 'Device name', value: 'name'},
  {text: 'Type', value: 'type'}
];

const selectedDevices = computed(() => {
  return props.zone.devices;
});
</script>

<style scoped>
.v-list-item:not(:last-child) {
  border-bottom: 1px solid var(--v-neutral-lighten3);
}
</style>
