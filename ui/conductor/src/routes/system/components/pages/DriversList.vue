<template>
  <ServicesList :name="ServiceNames.Drivers"/>
</template>

<script>
import ServicesList from '@/routes/system/components/ServicesList.vue';
import {ServiceNames} from '@/api/ui/services';

export default {
  name: 'DriversList',
  components: {ServicesList},
  computed: {
    ServiceNames() {
      return ServiceNames;
    }
  }
};
</script>

<style scoped>

</style>
