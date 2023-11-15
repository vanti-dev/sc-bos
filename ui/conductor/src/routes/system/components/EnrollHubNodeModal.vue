<template>
  <v-dialog v-model="dialogState" class="elevation-0" width="auto" max-width="975px">
    <content-card v-if="dialogState">
      <v-row class="py-4 px-10 mb-2">
        <v-card-title>
          {{ props.certificateQuery.isQueried ? 'View Details' : 'Manage Component' }}
        </v-card-title>
        <v-spacer/>
        <v-btn class="mr-2 mt-3" icon @click="dialogState = false">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-row>
      <div>
        <component-input
            v-if="!readCertificates.length && !props.certificateQuery.isQueried"
            :address.sync="address"
            :dialog-state.sync="dialogState"
            :list-items="props.listItems"
            @inspectHubNodeAction="inspectHubNodeAction"
            @forgetHubNodeAction="forgetHubNodeAction"
            style="min-width: 450px;"/>
        <div v-if="readMetadata || readCertificates.length">
          <metadata-details v-if="readMetadata" :metadata="readMetadata"/>
          <certificate-details
              v-if="readCertificates"
              :address="address"
              :certificate-query="props.certificateQuery"
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
  certificateQuery: {
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

watch(() => props.certificateQuery, async (newValue) => {
  if (!newValue.address) {
    return;
  }

  await inspectHubNodeAction(newValue.address);
}, {immediate: true, deep: true});
</script>
