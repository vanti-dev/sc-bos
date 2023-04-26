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
            Automations: {{ automationTrackers[node.name]?.metadataTracker?.response?.totalCount }}
          </v-list-item>
        </v-list>
      </v-card-text>
    </v-card>
    <div/>
  </div>
</template>

<script setup>
import {ServiceNames} from '@/api/ui/services';
import {useErrorStore} from '@/components/ui-error/error';
import {useHubStore} from '@/stores/hub';
import {useServicesStore} from '@/stores/services';
import {computed, reactive, set} from 'vue';

const hubStore = useHubStore();
const servicesStore = useServicesStore();
const errorStore = useErrorStore();

const automationTrackers = reactive({});

const nodesList = computed(() => {
  return Object.values(hubStore.nodesList).map(node => {
    Promise.all([node.commsAddress, node.commsName])
        .then(([address, name]) => {
          console.debug('node', node);
          set(automationTrackers, node.name, servicesStore.getService(ServiceNames.Automations, address, name));
          errorStore.registerTracker(automationTrackers[node.name].metadataTracker);
          return servicesStore.refreshMetadata(ServiceNames.Automations, address, name);
        })
        .catch(e => {
          // should be caught by error store
        });
    return {
      ...node
    };
  });
});


</script>

<style scoped>

</style>
