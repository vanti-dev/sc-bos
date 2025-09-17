<template>
  <v-card elevation="0" tile>
    <v-card-title class="d-flex text-title-caps-large text-neutral-lighten-3">
      <span>On Off</span>
    </v-card-title>
    <v-list tile class="ma-0 pa-0">
      <v-list-item class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">State</v-list-item-title>
        <template #append>
          <v-list-item-subtitle>
            {{ stateStr }}
          </v-list-item-subtitle>
        </template>
      </v-list-item>
    </v-list>
    <v-card-actions class="px-4">
      <v-spacer/>
      <v-btn
          size="small"
          variant="tonal"
          @click="setOff()">
        Off
      </v-btn>
      <v-btn
          size="small"
          variant="tonal"
          @click="setOn()">
        On
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup>

import {useOccupancy} from '@/traits/occupancy/occupancy.js';
import {OnOff} from '@smart-core-os/sc-api-grpc-web/traits/on_off_pb.js';

const props = defineProps({
  value: {
    type: Object, // of OnOff.AsObject
    default: () => {
    }
  },
  loading: {
    type: Boolean,
    default: false
  },
  name: {
    type: String,
    required: true
  }
});
const emit = defineEmits([
  'updateOnOff' // of OnOff.AsObject | UpdateOnOffRequest.AsObject
]);

const {stateStr} = (() => props.value);

function setOff() {
  emit('updateOnOff', {name: props.name, onOff: {state: OnOff.State.OFF}});
}

function setOn() {
  emit('updateOnOff', {name: props.name, onOff: {state: OnOff.State.ON}});
}

</script>

<style scoped>
.v-list-item {
  min-height: auto;
}

.v-list-item__subtitle.occupied {
  color: rgb(var(--v-theme-success-lighten-1)) !important;
}

.v-list-item__subtitle.idle {
  color: rgb(var(--v-theme-info)) !important;
}

.v-list-item__subtitle.unoccupied {
  color: rgb(var(--v-theme-warning)) !important;
}
</style>