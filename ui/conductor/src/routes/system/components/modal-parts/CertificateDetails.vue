<template>
  <div class="d-flex flex-column align-center">
    <v-list width="100%" active-color="primary">
      <v-list-group>
        <template #activator="{props: _props, isOpen: _isOpen}">
          <v-list-item
              rounded="xl"
              :active="activeCertificate === 'root'"
              @click="setActiveCertificate(rootCertificate, 'root')">
            <template #prepend>
              <v-btn
                  v-bind="_props"
                  size="small"
                  variant="text"
                  :icon="_isOpen ? 'mdi-chevron-down' : 'mdi-chevron-right'"
                  class="mr-8">
                <v-icon size="22"/>
              </v-btn>
            </template>
            <v-list-item-title>
              {{ rootCertificate?.subject?.commonName }}
            </v-list-item-title>
            <template #append>
              <span class="font-weight-bold">Valid:</span>
              <v-icon end :color="checkValidity(rootCertificate?.validityPeriod?.to).color">
                {{ checkValidity(rootCertificate?.validityPeriod?.to).icon }}
              </v-icon>
            </template>
          </v-list-item>
        </template>
        <v-list-item
            v-for="(intermediateValue, intermediateKey) in intermediateCertificates"
            :key="intermediateKey"
            rounded="xl"
            class="mt-2"
            :active="activeCertificate === intermediateValue?.subject?.commonName"
            @click="setActiveCertificate(intermediateValue, intermediateValue?.subject?.commonName)">
          <template #prepend>
            <v-icon
                start
                :icon="activeCertificate === intermediateValue?.subject?.commonName ?
                  'mdi-chevron-down' : 'mdi-chevron-right'"/>
          </template>
          <v-list-item-title>
            {{ intermediateValue?.subject?.commonName }}
          </v-list-item-title>
          <template #append>
            <span class="font-weight-bold">Valid:</span>
            <v-icon end :color="checkValidity(intermediateValue?.validityPeriod?.to).color">
              {{ checkValidity(intermediateValue?.validityPeriod?.to).icon }}
            </v-icon>
          </template>
        </v-list-item>
      </v-list-group>
    </v-list>
    <div class="mb-2 pt-1 pb-4 px-6">
      <v-divider v-if="displayedCertificate" class="mt-5"/>
      <div v-if="displayedCertificate" class="pt-2 pb-4">
        <template v-for="(value, key) in displayedCertificate" :key="key">
          <v-list-item
              v-if="value"
              class="ma-0 pa-0 mb-n4">
            <div class="d-flex flex-row flex-nowrap align-start">
              <v-col cols="align-self" class="ma-0 pa-0 mr-n4">
                <v-list-item-title class="text-capitalize font-weight-bold ma-0 pa-0">
                  {{ camelToSentence(formatFingerprint(key)) }}:
                </v-list-item-title>
              </v-col>
              <v-col cols="10" class="ma-0 pa-0 pl-6">
                <v-list-item-subtitle
                    v-if="typeof value !== 'object'"
                    class="ma-0 pa-0 text-wrap">
                  {{ value }}
                </v-list-item-subtitle>
                <div v-else class="d-flex flex-column">
                  <template v-for="(subValue, subKey) in value" :key="subKey">
                    <v-list-item v-if="subValue" class="ma-0 pa-0 mt-n3 mb-n2">
                      <div class="d-flex flex-row pb-4">
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
                          <v-list-item-subtitle class="ma-0 pa-0 text-wrap">
                            {{ subValue }}
                          </v-list-item-subtitle>
                        </v-col>
                      </div>
                    </v-list-item>
                  </template>
                </div>
              </v-col>
            </div>
          </v-list-item>
        </template>
      </div>
    </div>

    <v-row class="mt-4" v-if="!props.nodeQuery.isQueried && !props.nodeQuery.isToForget">
      <v-card-text>
        Please check the information provided by the node is as you expect before clicking Enroll.
      </v-card-text>
    </v-row>
    <v-card-actions
        v-if="!props.nodeQuery.isQueried && !props.nodeQuery.isToForget"
        class="d-flex flex-row justify-space-between mt-4">
      <v-btn
          class="ml-n8 pl-4 pr-6"
          color="neutral-lighten-4"
          variant="text"
          @click="emits('resetCertificates')">
        <v-icon>mdi-chevron-left</v-icon>
        Back
      </v-btn>
      <v-btn
          class="px-4"
          color="primary"
          variant="text"
          @click="confirmEnroll">
        Enroll
      </v-btn>
    </v-card-actions>
  </div>
</template>

<script setup>
import {camelToSentence} from '@/util/string';
import {computed, ref, watchEffect} from 'vue';

const emits = defineEmits(['resetCertificates', 'enrollHubNodeAction']);
const props = defineProps({
  nodeQuery: {
    type: Object,
    default: () => ({})
  },
  readCertificates: {
    type: Array,
    default: () => []
  }
});
const _address = defineModel('address', {
  type: String,
  default: null
});

const activeCertificate = ref(null);
const displayedCertificate = ref(null);

// if the subject equals to the issuer, then it's the root certificate
const rootCertificate = computed(() => {
  const root = props.readCertificates.find((certificate) => {
    return certificate.subject.commonName === certificate.issuer.commonName;
  });
  if (root) {
    return root;
  }
  return null;
});

// if the subject equals to the issuer of the root certificate, then it's the intermediate certificate
const intermediateCertificates = computed(() => {
  const intermediates = props.readCertificates.filter((certificate) => {
    return certificate.subject.commonName !== certificate.issuer.commonName &&
        certificate.issuer.commonName === rootCertificate.value.subject.commonName;
  });
  if (intermediates.length > 0) {
    return intermediates;
  }
  return null;
});

const checkValidity = (date) => {
  const currentDate = new Date();
  const expirationDate = new Date(date);
  if (expirationDate < currentDate) {
    return {
      icon: 'mdi-close-circle',
      color: 'error'
    };
  }
  return {
    icon: 'mdi-check-circle',
    color: 'success-lighten-2'
  };
};

const formatFingerprint = (fingerprint) => {
  if (fingerprint.includes('sha1')) {
    return fingerprint.replace('sha1', 'sha-1');
  } else if (fingerprint.includes('sha256')) {
    return fingerprint.replace('sha256', 'sha-256');
  }
  return fingerprint;
};

const setActiveCertificate = (certificate, certificateName) => {
  activeCertificate.value = certificateName;
  displayedCertificate.value = certificate;
};

const resetActiveCertificate = () => {
  displayedCertificate.value = null;
  activeCertificate.value = null;
};

// watch the effect of rootCertificate and set it if no active certificate is set
watchEffect(() => {
  if (intermediateCertificates.value && activeCertificate.value === null) {
    setActiveCertificate(intermediateCertificates.value[0],
        intermediateCertificates.value[0].subject.commonName
    );
  } else if (rootCertificate.value && activeCertificate.value === null) {
    setActiveCertificate(rootCertificate.value, 'root');
  }
});

const confirmEnroll = () => {
  emits('enrollHubNodeAction', _address.value);
  emits('resetCertificates');
  resetActiveCertificate();
  _address.value = null;
};
</script>
