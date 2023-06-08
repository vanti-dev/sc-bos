<template>
  <content-card class="mb-8 d-flex flex-column px-6 pt-md-6">
    <v-row class="px-6 mb-6 pt-1">
      <h4 class="text-h4 pb-4 pb-lg-0 pt-0 pt-lg-4">Occupancy</h4>
      <v-spacer/>
      <v-col cols="3" class="d-flex flex-row flex-nowrap mx-0 px-0 justify-space-between" style="max-width: 160px;">
        <WithOccupancy v-slot="{ resource }" :name="props.name">
          <div class="text-h2 d-flex flex-row flex-nowrap align-end">
            {{ collectOccupancyData(resource.value) }}
            <span class="text-body-1 neutral--text text--lighten-5 pb-1 ml-1">/ {{ maxOccupancy }}</span>
          </div>
        </WithOccupancy>
        <div class="text-h2 d-flex flex-row flex-nowrap align-end">
          {{ occupancyPercentageDisplay }}
          <span class="text-body-1 neutral--text text--lighten-9 pb-1 ml-1">%</span>
        </div>
      </v-col>
    </v-row>
    <OccupancyGraph class="flex-grow-1 d-none d-md-block" width="100%" :name="props.name"/>
  </content-card>
</template>

<script setup>
import {computed, ref} from 'vue';
import ContentCard from '@/components/ContentCard.vue';
import WithOccupancy from '@/routes/devices/components/renderless/WithOccupancy.vue';
import OccupancyGraph from '@/routes/ops/components/OccupancyGraph.vue';

const props = defineProps({
  name: {
    type: String,
    default: 'building'
  }
});

const peopleCount = ref(0);
const maxOccupancy = 1234;

const occupancyPercentage = computed(() => (peopleCount.value / maxOccupancy) * 100);
const occupancyPercentageDisplay = computed(() =>
  occupancyPercentage.value > 0 ? occupancyPercentage.value.toFixed(1) : occupancyPercentage.value.toFixed(0)
);

/**
 *
 * @param {*} value
 * @return {number}
 */
function collectOccupancyData(value) {
  if (value) peopleCount.value = value.peopleCount;
  return peopleCount.value;
}
</script>
