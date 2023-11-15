<template>
  <div>
    <v-row class="ml-0 pl-0 my-0">
      <h3 class="text-h3 pt-2 pb-6">Components</h3>
      <v-spacer/>
      <v-btn color="neutral" @click="showModal = true">Manage Component</v-btn>
    </v-row>

    <div class="d-flex flex-wrap ml-n2">
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
            <v-tooltip top>
              <template #activator="{ on }">
                <v-btn class="ml-auto mr-1 mt-1" icon v-on="on" @click="onShowCertificates(node.address)">
                  <v-icon class="pt-1" size="24">mdi-certificate-outline</v-icon>
                </v-btn>
              </template>
              <span>Component Details</span>
            </v-tooltip>
          </v-chip-group>
        </v-card-text>
      </v-card>
      <div/>
    </div>

    <EnrollHubNodeModal
        :show-modal.sync="showModal"
        :certificate-query.sync="certificateQuery"
        :list-items="nodesList"/>
  </div>
</template>

<script setup>
import {ref, watch} from 'vue';
import StatusAlert from '@/components/StatusAlert.vue';
import useSystemComponents from '@/composables/useSystemComponents';
import EnrollHubNodeModal from '@/routes/system/components/EnrollHubNodeModal.vue';

const showModal = ref(false);
const {
  nodeDetails,
  nodesList,
  isProxy,
  isHub
} = useSystemComponents();

const certificateQuery = ref({
  address: null,
  isQueried: false
});

const onShowCertificates = (address) => {
  certificateQuery.value.address = address;
  certificateQuery.value.isQueried = true;
  showModal.value = true;
};

watch(showModal, (newModal) => {
  if (newModal === false) {
    certificateQuery.value.address = null;
    certificateQuery.value.isQueried = false;
  }
}, {immediate: true, deep: true, flush: 'sync'});
</script>

<style scoped>

</style>
