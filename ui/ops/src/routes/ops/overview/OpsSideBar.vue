<template>
  <component :is="component"/>
</template>

<script setup>
import NotificationSideBar from '@/routes/ops/notifications/NotificationSideBar.vue';
import {useSidebarStore} from '@/stores/sidebar.js';
import {computed} from 'vue';

const sidebar = useSidebarStore();
// for backwards compatibility, default to the notification sidebar.
const component = computed(() => sidebar.component ?? NotificationSideBar);
</script>
