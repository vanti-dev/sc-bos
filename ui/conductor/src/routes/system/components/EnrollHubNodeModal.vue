<template>
  <v-dialog v-model="dialogState" class="elevation-0" width="auto" max-width="975px">
    <content-card>
      <v-row class="pa-4 mb-2">
        <v-card-title>
          {{ !readCertificates.length ? 'Manage Component' : 'Certificates' }}
        </v-card-title>
        <v-spacer/>
        <v-btn class="mr-2 mt-3" icon @click="dialogState = false">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-row>
      <component-input
          v-if="!readCertificates.length"
          :address.sync="address"
          :dialog-state.sync="dialogState"
          :list-items="props.listItems"
          @inspectHubNodeAction="inspectHubNodeAction"
          @forgetHubNodeAction="forgetHubNodeAction"
          style="min-width: 450px;"/>

      <certificate-details
          v-if="readCertificates.length > 0"
          :address="address"
          :read-certificates="readCertificates"
          @enrollHubNodeAction="enrollHubNodeAction"
          @resetCertificates="resetCertificates"/>
    </content-card>
  </v-dialog>
</template>

<script setup>
import {computed, ref} from 'vue';
import useSystemComponents from '@/composables/useSystemComponents';

import ContentCard from '@/components/ContentCard.vue';
import ComponentInput from '@/routes/system/components/modal-parts/ComponentInput.vue';
import CertificateDetails from '@/routes/system/components/modal-parts/CertificateDetails.vue';

const emits = defineEmits(['update:showModal']);
const props = defineProps({
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
  resetCertificates
} = useSystemComponents();
const address = ref(null);

const dialogState = computed({
  get() {
    return props.showModal;
  },
  set(value) {
    if (value === false) {
      address.value = null; // defaults to null when modal is closed
      resetCertificates(); // reset certificates when modal is closed
    }
    emits('update:showModal', value);
  }
});
</script>
