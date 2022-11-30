<template>
  <v-text-field
    v-model="search"
    prepend-inner-icon="mdi-magnify"
    label="Search ID, Location "
    single-line
    hide-details
    width="100%"
    outlined
    class="textarea sectionText"
  />
</template>

<script setup>
import { useLightingStore } from "@/stores/operate/lighting.js";
import { storeToRefs } from "pinia";

const store = useLightingStore();
const { search } = storeToRefs(store);
</script>
