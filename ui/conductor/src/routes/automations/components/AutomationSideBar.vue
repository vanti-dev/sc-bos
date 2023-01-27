<template>
  <side-bar>
    <edit-config-card/>
  </side-bar>
</template>

<script>
import SideBar from '@/components/SideBar.vue';
import EditConfigCard from '@/routes/automations/components/service-cards/EditConfigCard.vue';

export default {
  name: 'AutomationSideBar',
  components: {EditConfigCard, SideBar}
};
</script>

<style scoped>

</style>
