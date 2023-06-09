<template>
  <content-card class="mb-8 d-flex flex-column px-6 pt-md-6">
    <v-row class="px-6 mb-6 pt-1">
      <h4 class="text-h4 pb-4 pb-lg-0 pt-0 pt-lg-4">Occupancy</h4>
      <v-spacer/>
      <v-col cols="3" class="d-flex flex-row flex-nowrap mx-0 px-0 justify-space-between" style="max-width: 150px;">
        <WithOccupancy
            v-slot="{ resource }"
            :name="props.name">
          <OccupancyPeopleCount :people-count="resource?.value?.peopleCount"/>
        </WithOccupancy>
      </v-col>
    </v-row>
    <OccupancyGraph class="flex-grow-1 d-none d-md-block" width="100%" :name="props.name"/>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import WithOccupancy from '@/routes/devices/components/renderless/WithOccupancy.vue';
import OccupancyGraph from '@/routes/ops/components/OccupancyGraph.vue';
import OccupancyPeopleCount from './OccupancyPeopleCount.vue';

const props = defineProps({
  name: {
    type: String,
    default: 'building'
  }
});
</script>
