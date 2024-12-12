<template>
  <div>
    <v-row class="ml-0 pl-0 my-0">
      <h3 class="text-h3 pt-2 pb-6">Zones</h3>
    </v-row>
    <services-list :name="ServiceNames.Zones"/>
  </div>
</template>

<script setup>
import {ServiceNames} from '@/api/ui/services.js';
import ServicesList from '@/routes/system/components/ServicesList.vue';
</script>

<style scoped>

</style>
