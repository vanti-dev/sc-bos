<template>
  <div class="d-flex flex-wrap">
    <v-card v-for="node in nodesList" :key="node.name" width="300px" class="ma-2">
      <v-card-title class="text-body-large font-weight-bold">{{ node.name }}</v-card-title>
      <v-card-subtitle v-if="node.description !== ''">{{ node.description }}</v-card-subtitle>
      <v-card-text>
        <v-list dense>
          <v-list-item
              class="pa-0"
              style="min-height: 20px">
            {{ node.address }}
          </v-list-item>
          <v-list-item
              :class="[{'red--text': trackers.metadataTracker.error}, 'pa-0 ma-0']"
              style="min-height: 20px"
              v-for="(trackers, service) in nodeDetails[node.name]"
              :key="service">
            <span class="mr-1">{{ service }}: {{ trackers.metadataTracker?.response?.totalCount }}</span>
            <StatusAlert :resource="trackers.metadataTracker?.error"/>
          </v-list-item>
        </v-list>
        <v-chip-group>
          <v-chip v-if="isProxy(node.name)" color="accent" small>gateway</v-chip>
          <v-chip v-if="isHub(node.name) && !isProxy(node.name)" color="primary" small>hub</v-chip>
        </v-chip-group>
      </v-card-text>
    </v-card>
    <div/>
  </div>
</template>

<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import useSystemComponents from '@/composables/useSystemComponents';

const {
  nodeDetails,
  nodesList,
  isProxy,
  isHub
} = useSystemComponents();
</script>

<style scoped>

</style>
