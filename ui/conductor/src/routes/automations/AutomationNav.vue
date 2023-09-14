<template>
  <v-list class="pa-0" dense nav>
    <v-list-item :disabled="accessLevel('/automations/all').blockedAccess" to="/automations/all">
      <v-list-item-icon>
        <v-icon>mdi-view-list</v-icon>
      </v-list-item-icon>
      <v-list-item-content class="text-capitalize">All</v-list-item-content>
    </v-list-item>
    <v-list-item
        v-for="automation of automationTypeList"
        :key="automation.type"
        :to="'/automations/' + automation.type"
        class="my-2"
        :disabled="accessLevel('/automations/' + automation.type).blockedAccess">
      <v-list-item-icon>
        <v-icon v-if="icon.hasOwnProperty(automation.type)">{{ icon[automation.type] ?? defaultIcon }}</v-icon>
      </v-list-item-icon>
      <v-list-item-content class="text-capitalize text-truncate">
        {{ formatNaming(automation.type) }}
      </v-list-item-content>
    </v-list-item>
  </v-list>
</template>

<script setup>
import {ServiceNames} from '@/api/ui/services';
import {usePageStore} from '@/stores/page';
import {useServicesStore} from '@/stores/services';
import {storeToRefs} from 'pinia';
import {computed, ref, watch} from 'vue';
import useAuthSetup from '@/composables/useAuthSetup';

const {accessLevel} = useAuthSetup();

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
  list.sort();
  return list;
});

// map of icons to use for different automation sections
const icon = ref({
  bms: 'mdi-office-building-cog-outline',
  history: 'mdi-history',
  lightreport: 'mdi-file-chart-outline',
  lights: 'mdi-lightbulb',
  statusalerts: 'mdi-alert-circle-outline',
  udmi: 'mdi-transit-connection-variant'
});
const defaultIcon = 'mdi-auto-mode';

const acronyms = ['bms', 'udmi'];
const suffixes = ['report', 'reports', 'alert', 'alerts'];

/**
 * @param {string} name
 * @return {string}
 */
const formatNaming = (name) => {
  for (const word of suffixes) {
    if (name.endsWith(word)) {
      return name.substring(0, name.length - word.length) + ' ' + word;
    }
  }

  for (const acronym of acronyms) {
    if (name === acronym) {
      return name.toUpperCase();
    }
  }

  return name;
};

watch(
    sidebarNode,
    async () => {
      metadataTracker.value = serviceStore.getService(
          ServiceNames.Automations,
          await sidebarNode.value.commsAddress,
          await sidebarNode.value.commsName
      ).metadataTracker;
      await serviceStore.refreshMetadata(
          ServiceNames.Automations,
          await sidebarNode.value.commsAddress,
          await sidebarNode.value.commsName
      );
    },
    {immediate: true}
);
</script>

<style scoped>
:deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
