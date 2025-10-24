<template>
  <v-card>
    <v-card-title class="text-h4 my-2">
      <v-icon icon="mdi-wrench" start size="20"/>
      Health Checks
    </v-card-title>
    <v-card-text>
      <health-impact-totals/>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import HealthImpactTotals from "@/traits/health/HealthImpactTotals.vue";
</script>

<style scoped>

</style>