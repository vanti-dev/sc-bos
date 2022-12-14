<template>
  <v-text-field
      v-model="search"
      prepend-inner-icon="mdi-magnify"
      label="Search ID, Location "
      hide-details
      width="100%"
      filled
      class="textarea sectionText"
  />
</template>

<script setup>
import {useLightingStore} from "@/stores/devices/lighting.js";
import {storeToRefs} from "pinia";

const store = useLightingStore();
const { search } = storeToRefs(store);
</script>
