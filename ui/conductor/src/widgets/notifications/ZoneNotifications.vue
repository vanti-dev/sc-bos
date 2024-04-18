<template>
  <content-card class="zone-notifications">
    <notifications-table overview-page :zone="props.zone"/>
  </content-card>
</template>
<script setup>
import ContentCard from '@/components/ContentCard.vue';
import NotificationsTable from '@/routes/ops/notifications/NotificationsTable.vue';

const props = defineProps({
  zone: {
    type: String,
    default: ''
  }
});
</script>
