<template>
      <v-card color="#292F35">
      <v-card-text>
        Smart Core OS:
        <span class="font-weight-bold">{{ smartCoreStatus }}</span>
      </v-card-text>
    </v-card>
</template>


<script setup>
import { useLightingStore } from "@/stores/operate/lighting.js";
import { storeToRefs } from "pinia";

const store = useLightingStore();

const { smartCoreStatus } = storeToRefs(store);
</script>
