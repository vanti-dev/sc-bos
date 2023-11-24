<template>
  <v-dialog v-model="dialogState" class="elevation-0" width="auto" max-width="975px">
    <content-card v-if="dialogState">
      <v-row class="py-4 px-6 mb-2">
        <v-card-title>
          {{ modalTitle }}
        </v-card-title>
        <v-spacer/>
        <v-btn class="mr-2 mt-3" icon @click="dialogState = false">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-row>
      <div>
        <component-input
            v-if="showInput"
            :address.sync="address"
            :dialog-state.sync="dialogState"
            :list-items="props.listItems"
            :node-query="props.nodeQuery"
            @inspectHubNodeAction="inspectHubNodeAction"
            @forgetHubNodeAction="forgetHubNodeAction"
            style="min-width: 450px;"/>

        <!-- Node details -->
        <div v-if="showDetails">
          <metadata-details v-if="readMetadata" :metadata="readMetadata"/>
          <certificate-details
              v-if="readCertificates"
              :address.sync="address"
              :node-query="props.nodeQuery"
              :read-certificates="readCertificates"
              @enrollHubNodeAction="enrollHubNodeAction"
              @resetCertificates="resetCertificates"/>
        </div>
      </div>
    </content-card>
  </v-dialog>
</template>

<script setup>
import {computed, ref, watch} from 'vue';
import useSystemComponents from '@/composables/useSystemComponents';

import ContentCard from '@/components/ContentCard.vue';
import ComponentInput from '@/routes/system/components/modal-parts/ComponentInput.vue';
import CertificateDetails from '@/routes/system/components/modal-parts/CertificateDetails.vue';
import MetadataDetails from '@/routes/system/components/modal-parts/MetadataDetails.vue';

const emits = defineEmits(['update:showModal']);
const props = defineProps({
  nodeQuery: {
    type: Object,
    default: () => ({})
  },
  showModal: {
    type: Boolean,
    required: true
  },
  listItems: {
    type: Array,
    default: () => []
  }
});

const {
  enrollHubNodeAction,
  forgetHubNodeAction,
  inspectHubNodeAction,
  readCertificates,
  resetCertificates,
  readMetadata
} = useSystemComponents();
const address = ref(null);


const dialogState = computed({
  get() {
    return props.showModal;
  },
  set(value) {
    if (value === false) {
      // defaults to null when modal is closed
      address.value = null;
      resetCertificates(); // reset certificates when modal is closed
    }
    emits('update:showModal', value);
  }
});

const modalTitle = computed(() => {
  return props.nodeQuery.isQueried ?
      'View Details' : props.nodeQuery.isToForget ?
          'Forget node' : 'Enroll a new node';
});

const showInput = computed(() => {
  return readCertificates.value.length === 0 && !props.nodeQuery.isQueried;
});

const showDetails = computed(() => {
  return props.nodeQuery.isQueried || readCertificates.value.length > 0;
});

// Watch the nodeQuery object for changes
// If we have an address and it's not to forget, then we want to inspect the node
// If it is to forget, then we want to forget the node
watch(() => props.nodeQuery, async (newValue) => {
  if (!newValue.address || newValue.isToForget) {
    return;
  }

  await inspectHubNodeAction(newValue.address);
}, {immediate: true, deep: true});
</script>
