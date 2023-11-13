<template>
  <div class="d-flex flex-column align-center px-4">
    <v-expansion-panels class="ma-0 pa-0 mb-6" flat popout tile>
      <v-expansion-panel
          v-for="(item,i) in props.readCertificates"
          :key="i"
          class="rounded-lg mb-2"
          style="border: 1px solid grey;">
        <v-expansion-panel-header>
          <div class="pr-4">
            <div class="d-flex flex-row flex-wrap justify-space-between">
              <div class="">
                <span class="font-weight-bold mr-2">Subject:</span>
                <span>{{ item.subject.commonName }}</span>
              </div>
              <div class="">
                <span class="font-weight-bold mr-2">Serial:</span>
                <span>{{ item.serial }}</span>
              </div>
              <div class="">
                <span class="font-weight-bold mr-2">Valid:</span>
                <span>{{ item.validityPeriod }}</span>
              </div>
            </div>
          </div>
        </v-expansion-panel-header>
        <v-divider/>
        <v-expansion-panel-content class="mt-2">
          <template v-for="(value, key) in item">
            <v-list-item v-if="value" class="ma-0 pa-0 mb-n4" :key="key">
              <v-list-item-content class="d-flex flex-row flex-nowrap align-start">
                <v-col cols="align-self" class="ma-0 pa-0 mr-n4">
                  <v-list-item-title class="text-capitalize font-weight-bold ma-0 pa-0">
                    {{ camelToSentence(formatFingerprint(key)) }}:
                  </v-list-item-title>
                </v-col>
                <v-col cols="10" class="ma-0 pa-0 pl-6">
                  <v-list-item-subtitle v-if="typeof value !== 'object'" class="ma-0 pa-0">
                    {{ value }}
                  </v-list-item-subtitle>
                  <div v-else class="d-flex flex-column">
                    <template v-for="(subValue, subKey) in value">
                      <v-list-item v-if="subValue" class="ma-0 pa-0 mt-n3 mb-n2" :key="subKey">
                        <v-list-item-content class="d-flex flex-row pb-4">
                          <v-col cols="3" class="ma-0 pa-0">
                            <v-list-item-title
                                class="text-capitalize
                                text-body-2
                                font-weight-medium
                                ma-0
                                pa-0">
                              {{ camelToSentence(subKey) }}:
                            </v-list-item-title>
                          </v-col>
                          <v-col cols="9" class="ma-0 pa-0 pl-6">
                            <v-list-item-subtitle class="ma-0 pa-0">
                              {{ subValue }}
                            </v-list-item-subtitle>
                          </v-col>
                        </v-list-item-content>
                      </v-list-item>
                    </template>
                  </div>
                </v-col>
              </v-list-item-content>
            </v-list-item>
          </template>
        </v-expansion-panel-content>
        <v-divider v-if="i < props.readCertificates.length - 1"/>
      </v-expansion-panel>
    </v-expansion-panels>

    <v-card-actions v-if="!props.certificateQuery.isQueried" class="d-flex flex-row justify-space-between mt-4">
      <v-btn
          class="pl-4 pr-6"
          color="neutral lighten-4"
          text
          @click="emits('resetCertificates')">
        <v-icon>mdi-chevron-left</v-icon>
        Back
      </v-btn>
      <v-btn
          class="px-4"
          color="primary"
          text
          @click="confirmEnroll">
        Confirm
      </v-btn>
    </v-card-actions>
  </div>
</template>

<script setup>
import {camelToSentence} from '@/util/string';

const emits = defineEmits(['resetCertificates', 'enrollHubNodeAction']);
const props = defineProps({
  address: {
    type: String,
    default: null
  },
  certificateQuery: {
    type: Object,
    default: () => ({})
  },
  readCertificates: {
    type: Array,
    default: () => []
  }
});


const formatFingerprint = (fingerprint) => {
  if (fingerprint.includes('sha1')) {
    return fingerprint.replace('sha1', 'sha-1');
  } else if (fingerprint.includes('sha256')) {
    return fingerprint.replace('sha256', 'sha-256');
  }
  return fingerprint;
};

const confirmEnroll = () => {
  emits('enrollHubNodeAction', props.address);
  emits('resetCertificates');
};
</script>
