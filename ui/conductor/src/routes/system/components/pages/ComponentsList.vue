<template>
  <div>
    <v-row class="ml-0 pl-0 my-0">
      <h3 class="text-h3 pt-2 pb-6">Components</h3>
      <v-spacer/>
      <v-tooltip location="left">
        <template #activator="{ props }">
          <v-btn
              class="mr-4"
              color="primary"
              icon="mdi-plus"
              size="small"
              v-bind="props"
              @click="showModal = true">
            <v-icon size="24"/>
          </v-btn>
        </template>
        Enroll Node
      </v-tooltip>
    </v-row>

    <div class="d-flex flex-wrap ml-n2">
      <v-card v-for="node in nodesList" :key="node.name" width="300px" class="ma-2">
        <v-card-title
            class="text-body-large font-weight-bold d-flex align-center text-wrap"
            style="word-break: break-all">
          {{ node.name }}
          <v-spacer/>
          <v-menu min-width="175px">
            <template #activator="{ props }">
              <v-btn
                  icon="mdi-dots-vertical"
                  variant="text"
                  size="small"
                  v-bind="props">
                <v-icon size="24"/>
              </v-btn>
            </template>
            <v-list class="py-0">
              <v-list-item link @click="onShowCertificates(node.address)">
                <v-list-item-title>
                  View Certificate
                </v-list-item-title>
              </v-list-item>
              <v-list-item v-if="allowForget(node.name)" link @click="onForgetNode(node.address)">
                <v-list-item-title class="text-error">
                  Forget Node
                </v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
        </v-card-title>
        <v-card-subtitle v-if="node.description !== ''">{{ node.description }}</v-card-subtitle>

        <v-card-text>
          <v-list density="compact">
            <v-list-item
                class="pa-0"
                style="min-height: 20px">
              {{ node.address }}
            </v-list-item>
            <v-list-item
                :class="[{'text-red': trackers.metadataTracker.error}, 'pa-0 ma-0']"
                style="min-height: 20px"
                v-for="(trackers, service) in nodeDetails[node.name]"
                :key="service">
              <span class="mr-1">{{ service }}: {{ trackers.metadataTracker?.response?.totalCount }}</span>
              <status-alert :resource="trackers.metadataTracker?.error"/>
            </v-list-item>
          </v-list>
          <div>
            <v-chip v-if="isProxy(node.name)" color="accent" size="small" variant="flat">gateway</v-chip>
            <v-chip v-if="isHub(node.name) && !isProxy(node.name)" color="primary" size="small" variant="flat">
              hub
            </v-chip>
          </div>
        </v-card-text>
      </v-card>
    </div>

    <!-- Modal -->
    <enroll-hub-node-modal
        v-model:show-modal="showModal"
        v-model:node-query="nodeQuery"
        :list-items="nodesList"/>
  </div>
</template>

<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import useSystemComponents from '@/composables/useSystemComponents';
import EnrollHubNodeModal from '@/routes/system/components/EnrollHubNodeModal.vue';
import {ref, watch} from 'vue';

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
