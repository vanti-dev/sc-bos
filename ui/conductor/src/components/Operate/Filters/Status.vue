<template>
  <v-select
    v-model="status"
    :items="statuses"
    label="Status"
    class="sectionText"
    outlined
    hide-details="auto"
    style="width: 100px"
  ></v-select>
</template>

<script setup>
import { useLightingStore } from "@/stores/operate/lighting.js";
import { storeToRefs } from "pinia";

const store = useLightingStore();
const { status, statuses } = storeToRefs(store);
</script>
