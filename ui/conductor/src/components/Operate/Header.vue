<template>
  <div class="d-flex justify-space-between mb-2">
    <v-card color="transparent" flat elevation-0>
      <v-card-text>
        System:
        <span class="font-weight-bold">Lighting</span>
      </v-card-text>
    </v-card>
    <v-card color="#292F35">
      <v-card-text>
        Smart Core OS:
        <span class="font-weight-bold">{{ smartCoreStatus }}</span>
      </v-card-text>
    </v-card>
  </div>
</template>

<script setup>
import { useLightingStore } from "@/stores/operate/lighting.js";
import { storeToRefs } from "pinia";

const store = useLightingStore();

const { smartCoreStatus } = storeToRefs(store);
</script>
