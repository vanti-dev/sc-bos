<template>
  <content-card class="d-flex flex-column pa-8">
    <div class="d-flex flex-row mb-6 pl-4 pr-2 align-baseline">
      <h4 class="text-h4 ma-0">Occupancy</h4>
      <v-spacer/>
      <occupancy-people-count :people-count="peopleCount" :error="streamError" class="text-h2"/>
    </div>
    <occupancy-graph class="flex-grow-1" :source="historyName"/>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import {useOccupancy, usePullOccupancy} from '@/traits/occupancy/occupancy.js';
import OccupancyGraph from '@/widgets/occupancy/OccupancyGraph.vue';
import {computed} from 'vue';
import OccupancyPeopleCount from './OccupancyPeopleCount.vue';

const props = defineProps({
  source: {
    type: String,
    default: null
  },
  history: {
    type: String,
    default: null
  }
});

const sourceName = computed(() => props.source);
const historyName = computed(() => props.history ?? sourceName.value);

const {value, streamError} = usePullOccupancy(sourceName);
const {peopleCount} = useOccupancy(value);
</script>
