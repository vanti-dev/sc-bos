<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0" three-line>
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">UDMI Event</v-subheader>
      <v-list-item class="py-1" v-if="message.updateTime">
        <v-list-item-content class="py-0">
          <v-list-item-title class="text-body-small text-capitalize">Last updated</v-list-item-title>
          <v-list-item-subtitle class="text-capitalize">
            {{ Intl.DateTimeFormat('en-GB', {dateStyle:'short',timeStyle:'long'}).format(message.updateTime) }}
          </v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
      <v-list-item class="py-1" v-if="message.value">
        <v-list-item-content class="py-0">
          <v-list-item-title class="text-body-small text-capitalize">Topic</v-list-item-title>
          <v-list-item-subtitle class="text-capitalize">{{ message.value?.topic }}</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
      <v-list-item class="py-1" v-for="(value, key) in messagePayload" :key="key">
        <v-list-item-title class="text-body-small text-capitalize">{{ key }}</v-list-item-title>
        <v-list-item-subtitle class="text-capitalize">{{ value['present_value'] }}</v-list-item-subtitle>
      </v-list-item>
      <v-progress-linear color="primary" indeterminate :active="message.loading || message.value === null"/>
    </v-list>
  </v-card>
</template>

<script setup>

import {closeResource, newResourceValue} from '@/api/resource';
import {pullExportMessages} from '@/api/sc/traits/udmi';
import {useErrorStore} from '@/components/ui-error/error';
import {computed, onMounted, onUnmounted, reactive, watch} from 'vue';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  }
});

const message = reactive(newResourceValue());

const messagePayload = computed(() => {
  if (message.value === null) return {};
  return JSON.parse(message.value.payload);
});

// UI error handling
const errorStore = useErrorStore();
let unwatchMessageError;
onMounted(() => {
  unwatchMessageError = errorStore.registerValue(message);
});
onUnmounted(() => {
  if (unwatchMessageError) unwatchMessageError();
});

watch(() => props.name, async (name) => {
  // close existing stream if present
  closeResource(message);
  // create new
  if (name && name !== '') {
    pullExportMessages(name, message);
  }
}, {immediate: true});

onUnmounted(() => {
  closeResource(message);
});

</script>

<style scoped>
.v-list-item {
    min-height: auto;
}
</style>
