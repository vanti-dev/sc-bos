<template>
      <v-card class="card">
      <v-card-text>
        Smart Core OS:
        <span class="font-weight-bold">{{ smartCoreStatus }}</span>
      </v-card-text>
    </v-card>
</template>


<script setup>
import { useSmartCoreStore } from "@/stores/smartcore/smartcore.js";
import { storeToRefs } from "pinia";

const store = useSmartCoreStore();

const { smartCoreStatus } = storeToRefs(store);
</script>
