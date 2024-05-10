<template>
  <content-card class="zone-notifications">
    <notifications-table overview-page :force-query="forceQuery"/>
  </content-card>
</template>
<script setup>
import ContentCard from '@/components/ContentCard.vue';
import NotificationsTable from '@/routes/ops/notifications/NotificationsTable.vue';
import {computed} from 'vue';

const props = defineProps({
  forceQuery: {
    type: Object, /** @type {import('@sc-bos/ui-gen/proto/alerts_pb').Alert.Query.AsObject} */
    default: null
  },
  // legacy, use forceQuery instead. Here because it's driven directly from json config.
  zone: {
    type: String,
    default: ''
  }
});

const forceQuery = computed(() => {
  const q = props.forceQuery ?? {};
  if (props.zone) {
    q.zone = props.zone;
  }
  return q;
});
</script>
