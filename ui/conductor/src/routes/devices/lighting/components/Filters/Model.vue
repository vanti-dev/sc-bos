<template>
  <v-select
      v-model="model"
      :items="models"
      label="Type"
      class="sectionText"
      filled
      hide-details="auto"
      style="width: 100px"/>
</template>

<script setup>
import {useLightingStore} from '@/stores/devices/lighting.js';
import {storeToRefs} from 'pinia';

const store = useLightingStore();

const {model, models} = storeToRefs(store);
</script>
