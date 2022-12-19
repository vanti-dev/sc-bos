<template>
  <v-card class="rounded-lg">
    <v-card-text class="white--text">
      Smart Core OS:
      <span class="font-weight-bold text-uppercase" :class="getColor(smartCoreStatus)">{{ smartCoreStatus }}</span>
    </v-card-text>
  </v-card>
</template>

<script setup>
import { useSmartCoreStore } from "@/stores/smartcore/smartcore.js";
import { storeToRefs } from "pinia";

const store = useSmartCoreStore();

const { smartCoreStatus } = storeToRefs(store); 

const getColor = (status) => {
  switch (status) {
    case "Online":
      return "green--text";
    case "Offline":
      return "red--text";
    default:
      return "grey--text";
  }
};

</script>
