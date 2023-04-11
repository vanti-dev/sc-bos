<template>
  <content-card class="mb-8 d-flex flex-column px-6 pt-md-6">
    <h4 class="text-h4 order-lg-last pb-4 pb-lg-0 pt-0 pt-lg-4">Occupancy</h4>
    <v-progress-linear height="24" class="mb-3" v-model="occupancyPercentage"/>
    <div>
      <div class="text-h2" style="float:left">
        {{ occupancy }}
        <span class="text-caption neutral--text text--lighten-5">/{{ maxOccupancy }}</span>
      </div>
      <div class="text-h2" style="float:right">{{ occupancyPercentage.toFixed(0) }}%</div>
    </div>
  </content-card>
</template>

<script setup>

import {computed, onMounted, onUnmounted, reactive, watch} from 'vue';
import ContentCard from '@/components/ContentCard.vue';
import {closeResource, newResourceValue} from '@/api/resource';
import {pullOccupancy} from '@/api/sc/traits/occupancy';
import {useErrorStore} from '@/components/ui-error/error';

const props = defineProps({
  name: {
    type: String,
    default: 'building'
  }
});

const occupancyValue = reactive(newResourceValue());
// todo: where do we get this from?
const maxOccupancy = 1234;

const occupancy = computed(() => {
  return occupancyValue.value?.peopleCount ?? 0;
});

const occupancyPercentage = computed(() => occupancy.value/maxOccupancy*100);

watch(() => props.name, async (name) => {
  // close existing stream if present
  closeResource(occupancyValue);
  // create new stream
  if (name && name !== '') {
    pullOccupancy(name, occupancyValue);
  }
}, {immediate: true});

onUnmounted(() => {
  closeResource(occupancyValue);
});

// UI Error Handling
const errorStore = useErrorStore();
let unwatchErrors;
onMounted(() => {
  unwatchErrors = errorStore.registerValue(occupancyValue);
});
onUnmounted(() => {
  if (unwatchErrors) unwatchErrors();
});

</script>

<style scoped>

</style>
