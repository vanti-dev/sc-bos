<template>
  <div>
    <v-row class="ml-0 pl-0 my-0">
      <h3 class="text-h3 pt-2 pb-6">Components</h3>
      <v-spacer/>
      <v-tooltip left>
        <template #activator="{ on, attrs }">
          <v-btn
              class="mr-4"
              color="primary"
              fab
              small
              v-bind="attrs"
              v-on="on"
              @click="showModal = true">
            <v-icon>mdi-plus</v-icon>
          </v-btn>
        </template>
        Enroll Node
      </v-tooltip>
    </v-row>

    <div class="d-flex flex-wrap ml-n2">
      <v-card v-for="node in nodesList" :key="node.name" width="300px" class="ma-2">
        <div class="d-flex flex-row align-center pt-2 mb-n4">
          <v-card-title class="text-body-large font-weight-bold">{{ node.name }}</v-card-title>
          <v-card-subtitle v-if="node.description !== ''">{{ node.description }}</v-card-subtitle>
          <v-menu min-width="175px" nudge-bottom="10" nudge-right="10" offset-y>
            <template #activator="{ on, attrs }">
              <v-btn
                  class="ml-auto mr-3"
                  icon
                  v-bind="attrs"
                  v-on="on">
                <v-icon size="24">mdi-dots-vertical</v-icon>
              </v-btn>
            </template>
            <v-list class="py-0">
              <v-list-item link>
                <v-list-item-title @click="onShowCertificates(node.address)">
                  View Certificate
                </v-list-item-title>
              </v-list-item>
              <v-list-item v-if="allowForget(node.name)" link>
                <v-list-item-title class="error--text" @click="onForgetNode(node.address)">
                  Forget Node
                </v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
        </div>

        <!--        <v-card-title class="text-body-large font-weight-bold">{{ node.name }}</v-card-title>-->

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
    </div>

    <!-- Modal -->
    <EnrollHubNodeModal
        :show-modal.sync="showModal"
        :node-query.sync="nodeQuery"
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
  isHub,
  allowForget
} = useSystemComponents();

const nodeQuery = ref({
  address: null,
  isQueried: false,
  isToForget: false
});

const onShowCertificates = (address) => {
  nodeQuery.value.address = address;
  nodeQuery.value.isQueried = true;
  nodeQuery.value.isToForget = false;
  showModal.value = true;
};

const onForgetNode = (address) => {
  nodeQuery.value.address = address;
  nodeQuery.value.isQueried = false;
  nodeQuery.value.isToForget = true;
  showModal.value = true;
};

watch(showModal, (newModal) => {
  if (newModal === false) {
    nodeQuery.value.address = null;
    nodeQuery.value.isQueried = false;
    nodeQuery.value.isToForget = false;
  }
}, {immediate: true, deep: true, flush: 'sync'});
</script>

<style scoped>

</style>
