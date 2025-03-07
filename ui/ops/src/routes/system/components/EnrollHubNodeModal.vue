<template>
  <v-dialog v-model="dialogState" class="elevation-0" width="auto" max-width="975px">
    <content-card v-if="dialogState">
      <div class="d-flex align-baseline mb-6">
        <v-card-title>
          {{ modalTitle }}
        </v-card-title>
        <v-spacer/>
        <v-btn icon="mdi-close" variant="text" @click="dialogState = false"/>
      </div>
      <div>
        <component-input
            v-if="showInput"
            v-model:address="address"
            v-model:dialog-state="dialogState"
            :inspect-hub-node-value="inspectHubNodeValue"
            :list-items="props.listItems"
            :node-query="_nodeQuery"
            @inspect-hub-node-action="inspectHubNodeAction"
            @reset-inspect-hub-node-value="resetInspectHubNodeValue"
            @forget-hub-node-action="forgetHubNodeAction"
            style="min-width: 450px;"/>

        <!-- Node details -->
        <div v-if="showDetails">
          <metadata-details v-if="readMetadata" :metadata="readMetadata" class="px-4 mb-4"/>
          <certificate-details
              v-if="readCertificates"
              v-model:address="address"
              :node-query="_nodeQuery"
              :read-certificates="readCertificates"
              @enroll-hub-node-action="enrollHubNodeAction"
              @reset-certificates="resetCertificates"/>
        </div>
      </div>
    </content-card>
  </v-dialog>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import useSystemComponents from '@/composables/useSystemComponents';
import CertificateDetails from '@/routes/system/components/modal-parts/CertificateDetails.vue';
import ComponentInput from '@/routes/system/components/modal-parts/ComponentInput.vue';
import MetadataDetails from '@/routes/system/components/modal-parts/MetadataDetails.vue';
import {computed, ref, watch} from 'vue';

const props = defineProps({
  listItems: {
    type: Array,
    default: () => []
  }
});
const _showModal = defineModel('showModal', {
  type: Boolean,
  default: false
});
const _nodeQuery = defineModel('nodeQuery', {
  type: Object,
  default: () => ({})
});

const {
  enrollHubNodeAction,
  forgetHubNodeAction,
  inspectHubNodeAction,
  inspectHubNodeValue,
  resetInspectHubNodeValue,
  readCertificates,
  resetCertificates,
  readMetadata
} = useSystemComponents();
const address = ref(null);


const dialogState = computed({
  get() {
    return _showModal.value;
  },
  set(value) {
    if (value === false) {
      // defaults to null when modal is closed
      address.value = null;
      resetCertificates(); // reset certificates when modal is closed
    }
    _showModal.value = value;
  }
});

const modalTitle = computed(() => {
  if (_nodeQuery.value.isQueried) {
    return 'View Details';
  } else if (_nodeQuery.value.isToForget) {
    return 'Forget node';
  } else {
    return 'Enroll a new node';
  }
});


const showInput = computed(() => {
  return readCertificates.value.length === 0 && !_nodeQuery.value.isQueried;
});

const showDetails = computed(() => {
  return _nodeQuery.value.isQueried || readCertificates.value.length > 0;
});

// Watch the _nodeQuery object for changes
// If we have an address and it's not to forget, then we want to inspect the node
// If it is to forget, then we want to forget the node
watch(_nodeQuery, async (newValue) => {
  if (!newValue.address || newValue.isToForget) {
    return;
  }

  await inspectHubNodeAction(newValue.address);
}, {immediate: true, deep: true});
</script>
