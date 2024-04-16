<template>
  <div>
    <content-card
        v-if="props.item.widgets.showEnergyConsumption"
        class="mt-1 mb-7 pa-8 d-flex flex-column"
        style="min-height:385px;">
      <power-history-graph
          chart-title="Energy Consumption"
          class="flex-grow-1"
          :metered="props.item.widgets.showEnergyConsumption"/>
    </content-card>
    <ContentCard v-if="props.item.widgets.showNotifications" class="mt-1 pt-4 pl-0">
      <notifications overview-page :zone="props.item.widgets.showNotifications"/>
    </ContentCard>
  </div>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import Notifications from '@/routes/ops/notifications/Notifications.vue';
import PowerHistoryGraph from '@/widgets/power-history/PowerHistoryGraph.vue';

const props = defineProps({
  item: {
    type: Object,
    default: () => ({})
  }
});
</script>
