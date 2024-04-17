<template>
  <content-card class="d-flex flex-column pa-8">
    <energy-graph
        v-if="props.showChart"
        :demand="props.metered"
        :generated="props.generated"
        class="flex-grow-1"/>
    <power-total
        v-if="props.showIntensity"
        :demand="props.metered"
        :generated="props.generated"
        :occupancy="props.metered"
        class="mx-auto mt-8"/>
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
