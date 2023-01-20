<template>
  <side-bar/>
</template>

<script setup>

import SideBar from '@/components/SideBar.vue';

const props = defineProps({
  accountId: {
    type: String,
    default: ''
  }
});

</script>

<style scoped>

</style>
