<template>
  <content-card class="d-flex flex-column pa-8">
    <div class="d-flex flex-row mb-6 pl-4 pr-2 align-baseline">
      <h4 class="text-h4 ma-0">Occupancy</h4>
      <v-spacer/>
      <occupancy-people-count :people-count="peopleCount" :error="streamError" class="text-h2"/>
    </div>
    <occupancy-graph class="flex-grow-1" :source="props.name"/>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import {useOccupancy, usePullOccupancy} from '@/traits/occupancy/occupancy.js';
import OccupancyGraph from '@/widgets/occupancy/OccupancyGraph.vue';
import OccupancyPeopleCount from './OccupancyPeopleCount.vue';

const props = defineProps({
  name: {
    type: String,
    default: 'building'
  }
});

const {value, streamError} = usePullOccupancy(() => props.name);
const {peopleCount} = useOccupancy(value);
</script>
