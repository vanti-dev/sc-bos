<template>
  <div>
    <v-row class="ml-0 pl-0 my-0">
      <h3 class="text-h3 pt-2 pb-6">Drivers</h3>
    </v-row>
    <ServicesList :name="ServiceNames.Drivers"/>
  </div>
</template>

<script setup>
import {ServiceNames} from '@/api/ui/services';
import ServicesList from '@/routes/system/components/ServicesList.vue';
</script>

<style scoped>

</style>
