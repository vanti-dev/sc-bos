<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-list-subheader class="text-title-caps-large text-neutral-lighten-3">Emergency Lighting</v-list-subheader>
      <v-list-item class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">Battery Level</v-list-item-title>
        <v-list-item-subtitle class="text-capitalize">{{ battery }}%</v-list-item-subtitle>
      </v-list-item>
    </v-list>
    <v-progress-linear
        height="34"
        class="mx-4 my-2"
        :model-value="battery"
        bg-color="neutral-lighten-1"
        color="accent"/>
    <v-list>
      <v-list-subheader class="text-body-large font-weight-bold">Testing History</v-list-subheader>
      <v-list-item v-for="([date, state, textColor]) in testHistory" :key="date" class="pb-2">
        <v-list-item-title class="text-body-small">{{ date }}</v-list-item-title>
        <v-list-item-subtitle class="text-title-caps" :class="textColor">{{ state }}</v-list-item-subtitle>
      </v-list-item>
      <v-list-item>
        <v-list-item-action>
          <v-btn color="green" :disabled="blockActions" size="small">Test Now</v-btn>
        </v-list-item-action>
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script setup>

import {closeResource, newResourceValue} from '@/api/resource';
import useAuthSetup from '@/composables/useAuthSetup';
import {computed, onUnmounted, reactive, watch} from 'vue';

const {blockActions} = useAuthSetup();

const props = defineProps({
  name: {
    type: String,
    default: ''
  }
});

const emergencyLightValue = reactive(newResourceValue());

watch(() => props.name, async (name) => {
  // close existing stream if present
  closeResource(emergencyLightValue);
  // create new stream
  if (name && name !== '') {
    // todo: implement the API for this
  }
}, {immediate: true});

onUnmounted(() => {
  closeResource(emergencyLightValue);
});

const battery = computed(() => {
  if (emergencyLightValue && emergencyLightValue.value) {
    return (emergencyLightValue.value.levelPercent * 100).toFixed(1);
  }
  return '-';
});

const testHistory = computed(() => {
  return [
    ['28.09.22', 'Pass', 'success--text text--lighten-3'],
    ['21.09.22', 'Pass', 'success--text text--lighten-3'],
    ['14.09.22', 'Fail', 'error--text text--lighten-1'],
    ['07.09.22', 'Pass', 'success--text text--lighten-3'],
    ['31.08.22', 'Pass', 'success--text text--lighten-3']
  ];
});

</script>

<style scoped>
.v-list-item {
  min-height: auto;
}

.v-progress-linear {
  width: auto;
}
</style>
