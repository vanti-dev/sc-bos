<template>
  <div>
    <content-card
        v-if="props.item.widgets.showEnergyConsumption"
        class="mt-1 pb-0 mb-7"
        style="min-height:385px;">
      <EnergyGraph
          chart-title="Energy Consumption"
          classes="pt-1 pb-2 ml-3 mr-5"
          :metered="props.item.widgets.showEnergyConsumption"/>
    </content-card>
    <ContentCard v-if="props.item.widgets.showNotifications" class="mt-1 pt-4 pl-0">
      <notifications overview-page :zone="props.item.widgets.showNotifications"/>
    </ContentCard>
  </div>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import EnergyGraph from '@/routes/ops/overview/pages/widgets/energyAndDemand/EnergyGraph.vue';
import Notifications from '@/routes/ops/notifications/Notifications.vue';

const props = defineProps({
  item: {
    type: Object,
    default: () => ({})
  }
});
</script>
