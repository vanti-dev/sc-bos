<template>
  <v-row class="ml-0 mb-4 mt-1">
    <LightingStatus
        v-if="props.item.traits.showEmergencyLighting"
        :item="props.item"/>
  </v-row>
</template>

<script setup>
import LightingStatus from '@/routes/ops/overview/pages/widgets/StatusWidgets/LightingStatus.vue';

const props = defineProps({
  item: {
    type: Object,
    default: () => ({})
  }
});
</script>
