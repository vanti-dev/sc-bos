<template>
  <services-list :name="ServiceNames.Automations" :type="type"/>
</template>

<script setup>
import {ServiceNames, startService, stopService} from '@/api/ui/services';
import ServicesList from '@/routes/system/components/ServicesList.vue';

const props = defineProps({
  type: {
    type: String,
    default: ''
  }
});

</script>

<style lang="scss" scoped>
</style>
