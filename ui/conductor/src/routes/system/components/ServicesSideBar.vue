<template>
  <side-bar>
    <status-card/>
    <v-divider/>
    <edit-config-card/>
  </side-bar>
</template>

<script setup>
import StatusCard from '@/routes/system/components/service-cards/StatusCard.vue';
import EditConfigCard from '@/routes/system/components/service-cards/EditConfigCard.vue';
import SideBar from '@/components/SideBar.vue';
</script>

<style scoped>

</style>
