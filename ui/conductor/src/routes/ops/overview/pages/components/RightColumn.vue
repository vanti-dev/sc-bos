<template>
  <div>
    <presence-card
        v-if="props.item.traits.showOccupancy"
        class="mb-5"
        :name="props.item.traits.showOccupancy"/>

    <environmental-card
        v-if="props.item.traits.showEnvironment"
        class="mt-3"
        gauge-color="#ffc432"
        :name="props.item.traits.showEnvironment.indoor"
        :external-name="props.item.traits.showEnvironment.outdoor"/>

    <content-card
        v-if="props.item.traits.showEnergyConsumption"
        class="pb-0"
        style="min-height:385px;">
      <v-card-title class="text-h4 pl-4">Energy Consumption</v-card-title>
      <EnergyGraph
          classes="mt-n2 ml-n2 mr-1"
          color="#ffc432"
          color-middle="rgba(255, 196, 50, 0.35)"
          :hide-legends="true"
          :metered="props.item.traits.showEnergyConsumption"/>
    </content-card>
  </div>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';

import EnvironmentalCard from '@/routes/ops/overview/pages/widgets/environmental/EnvironmentalCard.vue';
import PresenceCard from '@/routes/ops/overview/pages/widgets/occupancy/PresenceCard.vue';
import EnergyGraph from '@/routes/ops/overview/pages/widgets/energyAndDemand/EnergyGraph.vue';

const props = defineProps({
  item: {
    type: Object,
    default: () => ({})
  }
});
</script>
