<template>
  <ServicesList :name="ServiceNames.Drivers"/>
</template>

<script setup>
import {ServiceNames} from '@/api/ui/services';
import ServicesList from '@/routes/system/components/ServicesList.vue';
</script>

<style scoped>

</style>
