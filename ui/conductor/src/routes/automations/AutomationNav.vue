<template>
  <v-list class="pa-0" dense nav>
    <v-list-item to="/automations/all">
      <v-list-item-icon>
        <v-icon>mdi-view-list</v-icon>
      </v-list-item-icon>
      <v-list-item-content class="text-capitalize">All</v-list-item-content>
    </v-list-item>
    <v-list-item
        v-for="automation of automationTypeList"
        :key="automation.type"
        :to="'/automations/'+automation.type">
      <v-list-item-icon>
        <v-icon v-if="icon.hasOwnProperty(automation.type)">{{ icon[automation.type] }}</v-icon>
      </v-list-item-icon>
      <v-list-item-content class="text-capitalize">{{ automation.type }}</v-list-item-content>
    </v-list-item>
  </v-list>
</template>

<script setup>
import {ServiceNames} from '@/api/ui/services';
import {usePageStore} from '@/stores/page';
import {useServicesStore} from '@/stores/services';
import {storeToRefs} from 'pinia';
import {computed, ref, watch} from 'vue';

const serviceStore = useServicesStore();
const pageStore = usePageStore();
const {sidebarNode} = storeToRefs(pageStore);

const metadataTracker = ref({});

// filter out automations that have no instances, and map to {type, number} obj
const automationTypeList = computed(() => {
  if (!metadataTracker.value.response) return [];
  const list = [];
  metadataTracker.value.response.typeCountsMap.forEach(([type, number]) => {
    if (number > 0) {
      list.push({type, number});
    }
  });
  return list;
});

// map of icons to use for different automation sections
const icon = ref({
  lights: 'mdi-lightbulb',
  history: 'mdi-history'
});

watch(sidebarNode, async () => {
  console.log('sidebarNode', sidebarNode);
  metadataTracker.value = serviceStore.getService(
      ServiceNames.Automations,
      await sidebarNode.value.commsAddress,
      await sidebarNode.value.commsName).metadataTracker;
  await serviceStore.refreshMetadata(
      ServiceNames.Automations,
      await sidebarNode.value.commsAddress,
      await sidebarNode.value.commsName);
},
{immediate: true});

</script>

<style scoped>
:deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
