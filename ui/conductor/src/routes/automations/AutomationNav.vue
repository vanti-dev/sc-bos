<template>
  <v-list class="pa-0" dense nav>
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
import {useServicesStore} from '@/stores/services';
import {computed, onMounted, ref} from 'vue';
import {ServiceNames} from '@/api/ui/services';

const serviceStore = useServicesStore();
const metadataTracker = ref(serviceStore.getService(ServiceNames.Automations).metadataTracker);

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
  lights: 'mdi-lightbulb'
});

onMounted(() => serviceStore.refreshMetadata(ServiceNames.Automations));

</script>

<style scoped>
:deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
