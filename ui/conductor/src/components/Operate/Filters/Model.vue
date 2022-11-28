<template>
  <v-select
    v-model="model"
    :items="models"
    label="Type"
    class=" sectionText" 
    outlined
    hide-details="auto"
    style="width: 100px"
  ></v-select>
</template>

<script setup>
import { useLightingStore } from "@/stores/operate/lighting.js";
import { storeToRefs } from "pinia";

const store = useLightingStore();

const { model, models } = storeToRefs(store);
</script>
