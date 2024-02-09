<template>
  <div>
    <!--    <TopRow-->
    <!--        v-if="true"-->
    <!--        class="pb-4 mb-7"-->
    <!--        :item="props.item"/>-->

    <ContentCard>
      <h4 class="text-h4 pt-0 pt-lg-4 pl-4">Notifications</h4>
      <notifications overview-page :zone="props.item?.widgets?.showNotification"/>
    </ContentCard>
  </div>
</template>

<script setup>
// import TopRow from '@/routes/ops/overview/pages/components/TopRow.vue';
import ContentCard from '@/components/ContentCard.vue';
import Notifications from '@/routes/ops/notifications/Notifications.vue';

const props = defineProps({
  item: {
    type: Object,
    default: () => ({})
  }
});
</script>
