<template>
  <side-bar>
    <template v-if="canEdit" #actions>
      <v-btn :to="editLink" icon="mdi-pencil" variant="plain" size="small"/>
    </template>
    <status-card/>
    <v-divider/>
    <edit-config-card/>
  </side-bar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';
import EditConfigCard from '@/routes/system/components/service-cards/EditConfigCard.vue';
import StatusCard from '@/routes/system/components/service-cards/StatusCard.vue';
import {useSidebarStore} from '@/stores/sidebar.js';
import {computed} from 'vue';
import {useRoute} from 'vue-router';

const sidebar = useSidebarStore();
const route = useRoute();

const canEdit = computed(() => {
  return Boolean(route.meta?.editRoutePrefix && sidebar.data?.service?.id);
});
const editLink = computed(() => {
  if (!canEdit.value) return undefined;
  if (sidebar.data?.nodeName) {
    return {
      name: `${route.meta.editRoutePrefix}-name-id`, params: {
        id: sidebar.data.service.id,
        name: sidebar.data.nodeName
      }
    };
  } else {
    return {
      name: `${route.meta.editRoutePrefix}-id`, params: {
        id: sidebar.data.service.id
      }
    };
  }
});
</script>

<style scoped>

</style>
