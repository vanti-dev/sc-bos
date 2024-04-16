<template>
  <content-card class="d-flex flex-column pa-8">
    <div class="d-flex flex-row mb-6 pl-4 pr-2 align-baseline">
      <h4 class="text-h4 ma-0">Occupancy</h4>
      <v-spacer/>
      <WithOccupancy
          v-slot="{ resource }"
          :name="props.name">
        <OccupancyPeopleCount :people-count="resource?.value?.peopleCount"/>
      </WithOccupancy>
    </div>
    <OccupancyGraph class="flex-grow-1" :name="props.name"/>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import OccupancyGraph from '@/routes/ops/overview/pages/widgets/occupancy/OccupancyGraph.vue';
import WithOccupancy from '@/traits/occupancy/WithOccupancy.vue';
import OccupancyPeopleCount from './OccupancyPeopleCount.vue';

const props = defineProps({
  name: {
    type: String,
    default: 'building'
  }
});
</script>
