<template>
  <v-dialog content-class="elevation-0" max-width="500" overlay-opacity="0.4" v-model="showStatus">
    <template #activator="{ on, attrs }">
      <v-btn class="py-1" style="text-align: center" text v-bind="attrs" v-on="on" @click="showStatus = !showStatus">
        <span class="text-title mr-1">Smart Core OS:</span>
        <span :class="`text-title-bold text-uppercase ${statusColor}`">{{ statusText }}</span>
      </v-btn>
    </template>

    <content-card class="popup__status">
      <v-card class="elevation-0" style="border: 2px solid var(--v-neutral-lighten2)" width="500px">
        <v-card-title class="text-subtitle-1 mb-0 pb-0 mt-n1 mb-2">Connection status</v-card-title>
        <v-row class="d-flex flex-row justify-center mx-4 my-3">
          <!-- Test Hub to Nodes connection -->
          <div class="d-flex flex-row">
            <v-chip class="neutral lighten-1" small>UI</v-chip>
            <v-divider class="my-auto mx-2" style="min-width: 10px"/>

            <!-- Test UI to Server connection -->
            <status-alert
                :color="uiToServerStatus.color"
                :icon="uiToServerStatus.icon"
                :resource="uiToServerStatus.resource"/>

            <v-divider class="my-auto mx-2" style="min-width: 10px"/>
            <v-chip class="neutral lighten-4" small>Server</v-chip>
            <v-divider class="my-auto mx-2" style="min-width: 10px"/>

            <!-- Test Server To Hub connection -->
            <status-alert
                :color="serverToHubStatus.color"
                :icon="serverToHubStatus.icon"
                :resource="serverToHubStatus.resource"/>

            <v-divider class="my-auto mx-2" style="min-width: 10px"/>
            <v-chip class="primary" small>Hub</v-chip>
            <v-divider class="my-auto mx-2" style="min-width: 10px"/>

            <!-- Test Hub to Nodes connection -->
            <v-chip class="neutral lighten-2" small>Nodes</v-chip>
          </div>
        </v-row>
      </v-card>
    </content-card>
  </v-dialog>
</template>

<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import {computed, ref} from 'vue';
import {formatErrorMessage} from '@/util/error';
import {statusCodeToString} from '@/components/ui-error/util';
import useSystemComponents from '@/composables/useSystemComponents';

import ContentCard from '@/components/ContentCard.vue';

const showStatus = ref(false);
const {nodesListCollection, hubNode, hubNodeValue} = useSystemComponents();

const uiToServerStatus = computed(() => {
  const error = nodesListCollection?.streamError;
  const response = nodesListCollection?.value;

  if (error && !response || error && response) {
    return {
      color: 'warning',
      icon: 'mdi-alert',
      resource: error
    };
  } else if (!error && response) {
    return {
      color: 'success',
      icon: 'mdi-check',
      resource: {
        status: {
          code: '200',
          message: 'Connected',
          name: 'UI'
        }
      }
    };
  }

  return {
    color: 'warning',
    icon: 'mdi-alert',
    resource: 'Unavailable'
  };
});

const serverToHubStatus = computed(() => {
  const error = nodesListCollection?.streamError;
  const response = nodesListCollection?.value;

  // If streamError is not null and value is also not null, then we have an error
  if (error && Object.keys(response).length > 0) {
    return {
      color: 'error',
      icon: 'mdi-close',
      resource: error
    };
  } else if (!error && response) {
    return {
      color: 'success',
      icon: 'mdi-check',
      resource: {
        status: {
          code: '200',
          message: 'Connected',
          name: hubNode.name
        }
      }
    };
  }

  return {
    color: 'warning',
    icon: 'mdi-alert',
    resource: 'Unavailable'
  };
});

const statusText = computed(() => {
  if (uiToServerStatus.value.status) {
    return 'Online';
  } else {
    return 'Offline';
  }
});

// Computed property for calculating the status color
const statusColor = computed(() => {
  if (uiToServerStatus.value.status && serverToHubStatus.value.status) {
    return 'success--text text--lighten-4';
  } else {
    return 'error--text';
  }
});
</script>

<style lang="scss" scoped>
.popup {
  height: 100%;

  &__status {
    position: absolute;
    width: 500px;
    height: auto;
    max-height: 600px;
    left: auto;
    right: 155px;
    top: 55px;

    &::after {
      content: '';
      position: absolute;
      top: -7px;
      right: 20px;
      width: 0;
      height: 0;
      border-style: solid;
      border-width: 0 15px 15px 0;
      border-color: transparent transparent var(--v-neutral-lighten2) transparent;
      transform: rotate(135deg);
    }
  }
}
</style>
