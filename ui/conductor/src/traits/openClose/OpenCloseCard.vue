<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0" two-line>
      <v-list-subheader class="text-title-caps-large text-neutral-lighten-3">Open / Closed</v-list-subheader>

      <v-list-item class="mt-n4 mb-4">
        <v-list-item-content class="d-flex flex-row flex-nowrap justify-space-between align-center py-0">
          <v-list-item-title class="text-body-small text-capitalize my-auto">Open Percentage</v-list-item-title>
          <v-list-item-subtitle class="text-subtitle-1 font-weight-medium text-wrap ml-2">
            {{ openStr }}
          </v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script setup>
import {useOpenClosePositions} from '@/traits/openClose/openClose.js';

const props = defineProps({
  value: {
    type: Object,
    default: () => {}
  }
});

const {openStr} = useOpenClosePositions(() => props.value);
</script>
