<template>
  <div class="d-flex flex-wrap">
    <v-card v-for="node in nodesList" :key="node.name" width="300px" class="ma-2">
      <v-card-title class="text-body-large font-weight-bold pb-0">{{ node.name }}</v-card-title>
      <v-card-subtitle v-if="node.description !== ''">{{ node.description }}</v-card-subtitle>
      <v-card-text>
        <v-list dense>
          <v-list-item class="pa-0">
            {{ node.address }}
          </v-list-item>
          <v-list-item class="pa-0">
            Automations: {{ node.automations?.metadataTracker?.response?.totalCount }}
          </v-list-item>
        </v-list>
      </v-card-text>
    </v-card>
    <div/>
  </div>
</template>

<script setup>
import {ServiceNames} from '@/api/ui/services';
import {useHubStore} from '@/stores/hub';
import {useServicesStore} from '@/stores/services';
import {computed, reactive} from 'vue';

const hubStore = useHubStore();
const servicesStore = useServicesStore();

const nodesList = computed(() => {
  return Object.values(hubStore.nodesList).map(node => {
    console.debug('node', node);
    const automations = reactive(
        servicesStore.getService(ServiceNames.Automations, node.commsAddress, node.commsName)
    );
    servicesStore.refreshMetadata(ServiceNames.Automations, node.commsAddress, node.commsName);
    return {
      ...node,
      automations
    };
  });
});


</script>

<style scoped>

</style>
