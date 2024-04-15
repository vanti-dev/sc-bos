<template>
  <content-card class="mb-8 d-flex flex-column pt-6 pl-6 pr-8">
    <energy-graph
        v-if="props.showChart"
        class="mt-8 mb-6"
        :height="props.showIntensity ? '275px' : '315px'"
        :generated="props.generated"
        :metered="props.metered"/>
    <power-total
        v-if="props.showIntensity"
        :metered="props.metered"
        :generated="props.generated"
        :occupancy="props.metered"
        class="mx-auto mt-10 mb-6"/>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import EnergyGraph from '@/widgets/power-history/PowerHistoryGraph.vue';
import PowerTotal from '@/widgets/power-total/PowerTotal.vue';

const props = defineProps({
  generated: {
    type: String,
    default: ''
  },
  metered: {
    type: String,
    default: 'building'
  },
  showChart: {
    type: Boolean,
    default: true
  },
  showIntensity: {
    type: Boolean,
    default: true
  }
});
</script>

<style scoped>

</style>
