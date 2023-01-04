<template>
  <div class="d-flex justify-space-between mb-2">
    <v-card color="transparent" flat elevation-0>
      <v-card-text>
        System:
        <span class="font-weight-bold">Lighting</span>
      </v-card-text>
    </v-card>
    <smart-core-status-card/>
  </div>
</template>
<script setup>
import SmartCoreStatusCard from '@/routes/ops/smartCoreStatusCard.vue';

</script>
