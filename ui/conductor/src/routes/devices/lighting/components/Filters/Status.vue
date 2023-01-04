<template>
  <v-select
      v-model="status"
      :items="statuses"
      label="Status"
      class="sectionText"
      filled
      hide-details="auto"
      style="width: 100px"/>
</template>

<script setup>
import {useLightingStore} from '@/stores/devices/lighting.js';
import {storeToRefs} from 'pinia';

const store = useLightingStore();
const {status, statuses} = storeToRefs(store);
</script>
