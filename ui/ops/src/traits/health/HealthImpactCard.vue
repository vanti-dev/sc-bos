<template>
  <v-card>
    <v-card-title class="text-h4 my-2">
      <v-icon icon="mdi-wrench" start size="20"/>
      Health Checks
    </v-card-title>
    <v-card-text>
      <v-alert type="error" color="grey-lighten-2" style="filter: blur(2px); opacity: 0.2">
        <template #prepend>
          <v-icon color="error" icon="mdi-alert-circle-outline" size="25"/>
        </template>
        <template #title>
          LIFE SAFETY: <span class="text-error ml-2">{{ '1 ISSUE DETECTED' }}</span>
        </template>
      </v-alert>
    </v-card-text>
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