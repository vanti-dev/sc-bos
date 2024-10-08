<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0" lines="three" density="compact">
      <v-list-subheader class="text-title-caps-large text-neutral-lighten-3">UDMI Event</v-list-subheader>
      <v-list-item class="py-1 mb-2" v-if="message.updateTime">
        <v-list-item-title class="text-body-small text-capitalize">Last updated</v-list-item-title>
        <v-list-item-subtitle class="text-capitalize text-body-1">
          {{ Intl.DateTimeFormat('en-GB', {dateStyle: 'short', timeStyle: 'long'}).format(message.updateTime) }}
        </v-list-item-subtitle>
      </v-list-item>
      <v-list-item class="py-1 mb-2" v-if="message.value">
        <v-list-item-title class="text-body-small text-capitalize">Topic</v-list-item-title>
        <v-list-item-subtitle class="text-capitalize">{{ message.value?.topic }}</v-list-item-subtitle>
      </v-list-item>
      <v-list-item class="py-1" v-for="(value, key) in messagePayload" :key="key" lines="one">
        <v-list-item-title class="text-body-small text-capitalize flex-fill">{{ key }}</v-list-item-title>
        <template #append>
          <v-list-item-subtitle class="text-capitalize text-end flex-fill text-body-1">
            {{ value['present_value'] ?? value }}
          </v-list-item-subtitle>
        </template>
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
    pullExportMessages({name, includeLast: true}, message);
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
