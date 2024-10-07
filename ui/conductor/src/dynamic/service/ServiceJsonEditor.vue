<template>
  <v-card color="transparent" :loading="loading">
    <v-card-title class="d-flex align-center">
      <span class="mr-auto">
        {{ readonly ? 'Viewing' : 'Editing' }} "{{ serviceID }}" of type {{ serviceType }}
      </span>
      <v-fade-transition>
        <v-btn v-if="refreshNeeded" v-bind="refreshAttrs"
               variant="text" class="mr-4"
               v-tooltip="refreshTooltip">
          Refresh
        </v-btn>
      </v-fade-transition>
      <span v-tooltip="saveTooltip">
        <v-btn v-bind="saveAttrs"/>
      </span>
    </v-card-title>
    <v-expand-transition>
      <v-alert v-if="alertVisible" v-bind="alertAttrs"/>
    </v-expand-transition>
    <v-card-text>
      <json-editor v-model="configModel" :read-only="readonly"/>
    </v-card-text>
  </v-card>
</template>

<script setup>
import JsonEditor from '@/components/JsonEditor.vue';
import {useServiceConfig} from '@/dynamic/service/service.js';

const {
  saveAttrs, saveTooltip,
  refreshAttrs, refreshTooltip, refreshNeeded,
  alertVisible, alertAttrs,
  loading, configModel, readonly,
  serviceID, serviceType, serviceDisabled
} = useServiceConfig((v) => {
  const unwanted = ['disabled', 'name', 'type'];
  return Object.fromEntries(Object.entries(v).filter(([key]) => !unwanted.includes(key)));
}, (v) => {
  v.name = serviceID.value;
  v.type = serviceType.value;
  v.disabled = serviceDisabled.value;
  return v;
});
</script>
