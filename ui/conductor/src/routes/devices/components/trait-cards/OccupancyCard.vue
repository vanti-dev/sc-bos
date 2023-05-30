<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Occupancy Sensor</v-subheader>
      <HotPoint :item="sidebarData" :item-key="props.name">
        <template #hotpoint>
          <WithOccupancy :name="props.name">
            <template #occupancy="{occupancyData}">
              <v-list-item class="py-1">
                <v-list-item-title class="text-body-small text-capitalize">State</v-list-item-title>
                <v-list-item-subtitle
                    :class="[
                      occupancyData.occupancyState.toLowerCase(),
                      'text-capitalize text-subtitle-2 py-1 font-weight-medium'
                    ]">
                  {{ occupancyData.occupancyState }}
                </v-list-item-subtitle>
              </v-list-item>
              <v-list-item class="py-1" v-if="occupancyData.occupantCount !== 0">
                <v-list-item-title class="text-body-small text-capitalize">Count</v-list-item-title>
                <v-list-item-subtitle class="text-capitalize">{{ occupancyData.occupantCount }}</v-list-item-subtitle>
              </v-list-item>
              <v-progress-linear color="primary" indeterminate :active="occupancyData.occupancyValue.loading"/>
            </template>
          </WithOccupancy>
        </template>
      </HotPoint>
    </v-list>
  </v-card>
</template>

<script setup>
import WithOccupancy from '../renderless/WithOccupancy.vue';
import HotPoint from '@/components/HotPoint.vue';
import {usePageStore} from '@/stores/page';

const pageStore = usePageStore();
const {sidebarData} = pageStore;

const props = defineProps({
  name: {
    type: String,
    default: ''
  }
});

</script>

<style scoped>
.v-list-item {
    min-height: auto;
}
.v-list-item__subtitle.occupied {
    color: var(--v-success-lighten1) !important;
}
.v-list-item__subtitle.idle {
    color: var(--v-info-base) !important;
}
.v-list-item__subtitle.unoccupied {
    color: var(--v-warning-base) !important;
}
</style>
